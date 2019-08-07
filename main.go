package main

import (
	"github.com/Shopify/terraform-provider-dyn/dyn"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dyn.Provider})
}
