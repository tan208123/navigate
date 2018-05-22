package plugin

import (
	"github.com/sirupsen/logrus"
	"github.com/tan208123/navigate/grpc/drivers/rke"
	"github.com/tan208123/navigate/grpc/types"
)

var (
	BuiltInDrivers = map[string]bool{
		//		"gke":    true,
		//		"aks":    true,
		"rke": true,
		//		"import": true,
	}
)

func Run(driverName string, addrChan chan string) (types.Driver, error) {
	var driver types.Driver
	switch driverName {
	case "rke":
		driver = rke.NewDriver()
	default:
		addrChan <- ""
	}
	if BuiltInDrivers[driverName] {
		go types.NewServer(driver, addrChan).Serve()
		return driver, nil
	}
	logrus.Fatal("driver not supported")
	return driver, nil
}
