package main

import (
	"github.com/davidfischer-ch/terraform-provider-awx/awx"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return awx.Provider()
		},
	})
}
