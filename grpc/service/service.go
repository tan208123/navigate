package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tan208123/navigate/grpc/cluster"
	"github.com/tan208123/navigate/grpc/plugin"
	"github.com/tan208123/navigate/grpc/types"
)

var (
	pluginAddress = map[string]string{}
	Drivers       = map[string]types.Driver{}
)

func Start() error {
	for driver := range plugin.BuiltInDrivers {
		logrus.Infof("Activating driver %s", driver)
		addr := make(chan string)
		rpcDriver, err := plugin.Run(driver, addr)
		if err != nil {
			return err
		}
		Drivers[driver] = rpcDriver
		listenAddr := <-addr
		pluginAddress[driver] = listenAddr
		logrus.Infof("Activating driver %s done", driver)
	}
	return nil
}

type EngineService interface {
	Create(ctx context.Context, name string, clusterSpec string) (string, string, string, error)
	Remove(ctx context.Context, name string, clusterSpec string) error
}

type engineService struct {
	serviceName string
}

func NewEngineService() EngineService {
	return &engineService{
		serviceName: "test",
	}
}

type controllerConfigGetter struct {
	driverName  string
	clusterSpec string
	clusterName string
}

func (c controllerConfigGetter) GetConfig() (types.DriverOptions, error) {
	driverOptions := types.DriverOptions{
		BoolOptions:        make(map[string]bool),
		StringOptions:      make(map[string]string),
		IntOptions:         make(map[string]int64),
		StringSliceOptions: make(map[string]*types.StringSlice),
	}

	switch c.driverName {
	case "rke":
		driverOptions.StringOptions["config-file-path"] = c.clusterSpec
	}

	// driverOptions.StringOptions["name"] = c.clusterName

	return driverOptions, nil
}

func (e *engineService) convertCluster(name string, spec string) (cluster.Cluster, error) {
	driverName := ""
	if spec != "" {
		driverName = "rke"
	}
	if driverName == "" {
		return cluster.Cluster{}, fmt.Errorf("no driver config found")
	}
	pluginAddr := pluginAddress[driverName]
	configGetter := controllerConfigGetter{
		driverName:  driverName,
		clusterSpec: spec,
		clusterName: name,
	}
	clusterPlugin, err := cluster.NewCluster(driverName, pluginAddr, name, configGetter)
	if err != nil {
		return cluster.Cluster{}, err
	}
	return *clusterPlugin, nil

}

// Create creates the stub for cluster manager to call
func (e *engineService) Create(ctx context.Context, name string, clusterSpec string) (string, string, string, error) {
	cls, err := e.convertCluster(name, clusterSpec)
	if err != nil {
		return "", "", "", err
	}
	if err := cls.Create(ctx); err != nil {
		return "", "", "", err
	}
	endpoint := cls.Endpoint
	if !strings.HasPrefix(endpoint, "https://") {
		endpoint = fmt.Sprintf("https://%s", cls.Endpoint)
	}
	return endpoint, cls.ServiceAccountToken, cls.RootCACert, nil
}

// Remove removes stub for cluster manager to call
func (e *engineService) Remove(ctx context.Context, name string, clusterSpec string) error {
	cls, err := e.convertCluster(name, clusterSpec)
	if err != nil {
		return err
	}
	return cls.Remove(ctx)
}
