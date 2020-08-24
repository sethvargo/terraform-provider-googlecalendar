package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/sethvargo/terraform-provider-googlecalendar/googlecalendar"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: googlecalendar.Provider,
	})
}
