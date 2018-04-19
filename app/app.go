package app

import (
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Config struct {
	KubeConfig string
	K8sMode    string
	Embedded   bool
}

func Run(config *Config) error {
	logrus.Infof("app run ... ")
	logrus.Infof("config is %v ", config)
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		logrus.Infof("ListenAndServe: %v ", err)
	}
	return nil
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}
