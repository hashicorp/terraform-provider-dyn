package main

import (
	"github.com/SplunkCloud/terraform-provider-dyn/dyn"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dyn.Provider})
}
