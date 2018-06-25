package config

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ManagementContext struct {
	RESTConfig *rest.Config
	K8sClient  kubernetes.Interface
}

func NewManagementContext(config *rest.Config) (*ManagementContext, error) {
	var err error

	context := &ManagementContext{
		RESTConfig: config,
	}

	context.K8sClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return context, err

}
