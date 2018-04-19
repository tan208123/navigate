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
	server := &http.Server{Addr: ":12345"}
	http.HandleFunc("/hello", HelloServer)
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

	<-ctx.Done()
	return ctx.Err()
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}
