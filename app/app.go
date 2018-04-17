package app

import (
	"fmt"
)

type Config struct {
	KubeConfig      string
}

func Run(config *Config) error {
	fmt.Printf("app run ... \n")
	fmt.Printf("config is %v \n", config)
	return nil
}