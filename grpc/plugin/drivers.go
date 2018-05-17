package plugin

import (
	"github.com/sirupsen/logrus"
)

var (
	BuiltInDrivers = map[string]bool {
		"gke":    true,
		"aks":    true,
		"rke":    true,
		"import": true,
	}
)

func Run(drivername string, addrChan chan string) (types.Driver, error) {
	
}