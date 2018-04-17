package main

import (	
	"os"
//	"log"
//	"net/http"

	"github.com/tan208123/navigate/app"
	"github.com/urfave/cli"
)

var (
	VERSION = "dev"
)

func main() {

	var config  app.Config

	app := cli.NewApp()
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			Usage:       "Kube config for accessing k8s cluster",
			EnvVar:      "KUBECONFIG",
			Destination: &config.KubeConfig,
		},
	}

	app.Action = func(c *cli.Context) error {
		/*
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
		*/
		return run(config)
	}

	app.Run(os.Args)
}

func run(config app.Config) error {
	// fmt.Printf("config is %v \n", config)
	return app.Run(&config) 
}