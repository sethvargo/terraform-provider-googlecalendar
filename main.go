package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sethvargo/terraform-provider-google-calendar/googlecalendar"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return googlecalendar.Provider()
		},
	})
}
