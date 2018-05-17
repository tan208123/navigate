package service

import (
	"github.com/tan208123/navigate/grpc/plugin"
	"github.com/sirupsen/logrus"
)

func Start() error{
    for driver := range plugin.BuiltInDrivers {
		logrus.Infof("Activating driver %s", driver)
		
	}
}