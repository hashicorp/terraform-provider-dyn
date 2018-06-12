package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/splunkcloud/terraform-provider-dyn/dyn"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dyn.Provider})
}
