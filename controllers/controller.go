package controllers

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tan208123/navigate/controllers/annotator"
	"github.com/tan208123/navigate/controllers/configgenerator"
	"github.com/tan208123/navigate/controllers/healthchecker"
	"github.com/tan208123/navigate/controllers/provisioner"
	client "github.com/tan208123/navigate/pkg/client/clientset/versioned"
	informers "github.com/tan208123/navigate/pkg/client/informers/externalversions"
	"github.com/tan208123/navigate/pkg/config"
)

func Register(management *config.ManagementContext) error {
	logrus.Infof("restconfig is %v ", management.RESTConfig)

	client, err := client.NewForConfig(management.RESTConfig)
	if err != nil {
		return err
	}
	clusterInformerFactory := informers.NewSharedInformerFactory(client, time.Second*30)
	provisioner.Register(client, clusterInformerFactory, management)
	configgenerator.Register(client, clusterInformerFactory)
	healthchecker.Register(client, clusterInformerFactory)
	annotator.Register(client, clusterInformerFactory)

	return nil
}
