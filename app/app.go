package app

import (
	"context"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tan208123/navigate/grpc/service"
	"google.golang.org/grpc/metadata"
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
	if err := service.Start(); err != nil {
		logrus.Infof("service started error", err)
	}
	server := &http.Server{Addr: ":12345"}
	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/create", ClusterCreate)
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
	logrus.Infof("context done ...")
	return ctx.Err()
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

//use http://<ip>:<port>/create?cluster_name
func ClusterCreate(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	clusterName := req.FormValue("cluster_name")
	rpcDriver := service.NewEngineService()
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{}))
	spec := "./config/cluster_rke.yml"
	rpcDriver.Create(ctx, clusterName, spec)
}
