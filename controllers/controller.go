package controllers

import (
	"time"

	"github.com/tan208123/navigate/controllers/provisioner"
	client "github.com/tan208123/navigate/pkg/client/clientset/versioned"
	informers "github.com/tan208123/navigate/pkg/client/informers/externalversions"
	"github.com/tan208123/navigate/pkg/config"
)

func Register(management *config.ManagementContext) error {

	client, err := client.NewForConfig(&management.RESTConfig)
	if err != nil {
		return err
	}
	clusterInformerFactory := informers.NewSharedInformerFactory(client, time.Second*30)
	provisioner.Register(client, clusterInformerFactory, management)

	return nil
}
