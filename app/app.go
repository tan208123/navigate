package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/tan208123/navigate/controllers"
	"github.com/tan208123/navigate/grpc/service"
	"github.com/tan208123/navigate/pkg/config"
	"k8s.io/client-go/rest"
)

type Config struct {
	KubeConfig string
	K8sMode    string
	Embedded   bool
	Debug      bool
}

func Run(ctx context.Context, kubeConfig rest.Config, cfg *Config) error {
	logrus.Infof("app run ... ")
	logrus.Infof("config is %v ", cfg)
	if err := service.Start(); err != nil {
		logrus.Infof("service started error", err)
	}
	server := &http.Server{Addr: ":12345"}
	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/create", ClusterCreate)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logrus.Infof("ListenAndServe: %v ", err)
		}
	}()

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	management, err := config.NewManagementContext(kubeConfig)

	// Create custom resource definitions
	if err := createCRDS(); err != nil {
		return err
	}

	// Register controllers
	if err := controllers.Register(management); err != nil {
		return err
	}

	<-ctx.Done()
	logrus.Infof("context done ...")
	return ctx.Err()
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

//use http://<ip>:<port>/create?cluster_name
func ClusterCreate(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	clusterName := req.FormValue("cluster_name")
	logrus.Infof("clusterName is %s", clusterName)
	//rpcDriver := service.NewEngineService()
	//ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{}))
	//spec := "./config/cluster_rke.yml"
	//rpcDriver.Create(ctx, clusterName, spec)
}

func createCRDS() error {
	logrus.Info("Creating CRDs...")
	cmdName := "kubectl"
	files, err := ioutil.ReadDir("./config/crd")
	if err != nil {
		return err
	}
	for _, file := range files {
		filePath := fmt.Sprintf("./config/crd/%s", file.Name())
		logrus.Infof("Creating crd for file %s", filePath)
		cmdArgs := []string{"apply", "-f", filePath}
		cmd := exec.Command(cmdName, cmdArgs...)
		var out bytes.Buffer
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create CRD [%s] %v %v", file.Name(), err, out.String())
		}
	}

	logrus.Info("Created CRDs")
	return nil
}
