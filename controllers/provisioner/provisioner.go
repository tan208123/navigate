package provisioner

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tan208123/navigate/grpc/service"
	types "github.com/tan208123/navigate/pkg/apis/clusterprovisioner/v1alpha1"
	client "github.com/tan208123/navigate/pkg/client/clientset/versioned"
	clusterclient "github.com/tan208123/navigate/pkg/client/clientset/versioned"
	informers "github.com/tan208123/navigate/pkg/client/informers/externalversions"
	listers "github.com/tan208123/navigate/pkg/client/listers/clusterprovisioner/v1alpha1"
	"github.com/tan208123/navigate/pkg/config"
	"github.com/tan208123/navigate/util"
	"google.golang.org/grpc/metadata"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

type Controller struct {
	clusterLister   listers.ClusterLister
	clusterInformer cache.SharedIndexInformer
	clusterClient   clusterclient.Interface
	syncQueue       *util.TaskQueue
	Driver          service.EngineService
}

func Register(management *config.ManagementContext) {

	clusterClient, err := client.NewForConfig(&management.RESTConfig)
	if err != nil {
		panic(err)
	}
	sampleInformerFactory := informers.NewSharedInformerFactory(clusterClient, time.Second*30)

	clusterInformer := sampleInformerFactory.Clusterprovisioner().V1alpha1().Clusters()

	controller := &Controller{
		clusterLister:   clusterInformer.Lister(),
		clusterInformer: clusterInformer.Informer(),
		clusterClient:   clusterClient,
		Driver:          service.NewEngineService(NewPersistentStore(management.K8sClient.CoreV1().Namespaces(), management.K8sClient.CoreV1())),
	}
	controller.syncQueue = util.NewTaskQueue(controller.sync)
	controller.clusterInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			controller.syncQueue.Enqueue(obj)
		},
		UpdateFunc: func(old, cur interface{}) {
			controller.syncQueue.Enqueue(cur)
		},
	})
	stop := make(chan struct{})
	go controller.syncQueue.Run(time.Second, stop)
	logrus.Infof("Registered %s controller", controller.getName())
}

func (c *Controller) getName() string {
	return "provisioner"
}

func (c *Controller) sync(key string) {
	cluster, err := c.clusterLister.Get(key)
	if err != nil {
		c.syncQueue.Requeue(key, err)
		return
	}

	if cluster.DeletionTimestamp != nil {
		err = c.handleClusterRemove(cluster)
	} else {
		err = c.handleClusterAdd(cluster)
	}
	if err != nil {
		c.syncQueue.Requeue(key, err)
		return
	}
}

func (c *Controller) handleClusterRemove(cluster *types.Cluster) error {
	logrus.Infof("Removing cluster %v", cluster.Name)
	if err := c.finalize(cluster, c.getName()); err != nil {
		return fmt.Errorf("error removing cluster %s %v", cluster.Name, err)
	} else {
		logrus.Infof("Successfully removed cluster %v", cluster.Name)
	}
	return nil
}

func (c *Controller) handleClusterAdd(cluster *types.Cluster) error {
	config, err := getConfigStr(cluster)
	if err != nil {
		return err
	}
	// Compare applied vs current config, and only run update when there are changes
	if config == cluster.Status.AppliedConfig {
		return nil
	}

	logrus.Infof("Cluster [%s] is updated; provisioning...", cluster.Name)
	// Add finalizer and other init fields
	if err := c.initialize(cluster, c.getName()); err != nil {
		return fmt.Errorf("error initializing cluster %s %v", cluster.Name, err)
	}

	// Provision the cluster
	_, err = types.ClusterConditionProvisioned.Do(cluster, func() (runtime.Object, error) {
		// this is the place where cluster provisioning backend logic is being invoked
		return cluster, c.provisionCluster(cluster)
	})

	if err != nil {
		return fmt.Errorf("error provisioning cluster %s %v", cluster.Name, err)
	}
	// Update cluster with applied spec
	if err := c.updateAppliedConfig(cluster, config); err != nil {
		return fmt.Errorf("error updating cluster %s %v", cluster.Name, err)
	}
	logrus.Infof("Successfully provisioned cluster %v", cluster.Name)
	return nil
}

func (c *Controller) removeCluster(cluster *types.Cluster) (err error) {
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{}))
	err = c.Driver.Remove(ctx, cluster.Name, cluster.Spec.ConfigPath)
	return err
}

func executeCommand(cmdName string, cmdArgs []string) (err error) {
	cmd := exec.Command(cmdName, cmdArgs...)
	var stdout io.ReadCloser
	stdout, err = cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error getting stdout from cmd '%v' %v", cmd, err)
	}
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("error starting cmd '%v' %v", cmd, err)
	}
	defer func() {
		err = cmd.Wait()
	}()
	printLogs(stdout)
	return err
}

func (c *Controller) provisionCluster(cluster *types.Cluster) (err error) {
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{}))
	apiEndpoint, serviceAccountToken, caCert, err := c.Driver.Create(ctx, cluster.Name, cluster.Spec.ConfigPath)
	//TODO
	//save apiEndpoint, serviceAccountToken, caCert to cluster.Status
	logrus.Infof("Cluster apiEndpoint is %s\n, serviceAccountToken is %s\n, caCert is %s", apiEndpoint, serviceAccountToken, caCert)
	return err
}

func getConfigStr(cluster *types.Cluster) (string, error) {
	b, err := ioutil.ReadFile(cluster.Spec.ConfigPath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func printLogs(r io.Reader) {
	buf := make([]byte, 80)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			fmt.Print(string(buf[0:n]))
		}
		if err != nil {
			break
		}
	}
}

func containsString(slice []string, item string) bool {
	for _, j := range slice {
		if j == item {
			return true
		}
	}
	return false
}

func (c *Controller) initialize(cluster *types.Cluster, finalizerKey string) error {
	//set finalizers
	metadata, err := meta.Accessor(cluster)
	if err != nil {
		return err
	}
	if containsString(metadata.GetFinalizers(), finalizerKey) {
		return nil
	}
	finalizers := metadata.GetFinalizers()
	finalizers = append(finalizers, finalizerKey)
	metadata.SetFinalizers(finalizers)
	for i := 0; i < 3; i++ {
		_, err = c.clusterClient.ClusterprovisionerV1alpha1().Clusters().Update(cluster)
		if err == nil {
			return err
		}
	}
	return nil
}

func (c *Controller) updateAppliedConfig(cluster *types.Cluster, config string) error {
	cluster.Status.AppliedConfig = config
	for i := 0; i < 3; i++ {
		_, err := c.clusterClient.ClusterprovisionerV1alpha1().Clusters().Update(cluster)
		if err == nil {
			return nil
		}
	}
	return nil
}

func (c *Controller) finalize(cluster *types.Cluster, finalizerKey string) error {
	toUpdate, err := c.clusterClient.ClusterprovisionerV1alpha1().Clusters().Get(cluster.Name, v1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	metadata, err := meta.Accessor(toUpdate)
	if err != nil {
		return err
	}
	// Check finalizer
	if metadata.GetDeletionTimestamp() == nil {
		// already deleted
		return nil
	}

	// already "finalized" by this controllerß
	if !containsString(metadata.GetFinalizers(), finalizerKey) {
		return nil
	}

	//run deletion hook - call cluster cleanup logic on the backend
	err = c.removeCluster(cluster)
	if err != nil {
		return err
	}
	// remove finalizer when/if the cleanup passed successfully
	var finalizers []string
	for _, finalizer := range metadata.GetFinalizers() {
		if finalizer == finalizerKey {
			continue
		}
		finalizers = append(finalizers, finalizer)
	}
	metadata.SetFinalizers(finalizers)

	for i := 0; i < 3; i++ {
		_, err = c.clusterClient.ClusterprovisionerV1alpha1().Clusters().Update(toUpdate)
		if err == nil {
			break
		}
	}

	return err
}
