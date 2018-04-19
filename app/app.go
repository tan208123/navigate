package app

import (
	"context"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Config struct {
	KubeConfig string
	K8sMode    string
	Embedded   bool
	Debug      bool
}

func Run(ctx context.Context, cfg *Config) error {
	logrus.Infof("app run ... ")
	logrus.Infof("config is %v ", cfg)
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		logrus.Infof("ListenAndServe: %v ", err)
	}
	<-ctx.Done()
	return ctx.Err()
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}
