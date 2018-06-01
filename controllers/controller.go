package controllers

import (
	"github.com/tan208123/navigate/controllers/provisioner"
	"github.com/tan208123/navigate/pkg/config"
)

func Register(management *config.ManagementContext) error {

	provisioner.Register(management)

	return nil
}
