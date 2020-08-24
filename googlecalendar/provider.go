// Package googlecalendar manages Google calendar events with Terraform.
package googlecalendar

import (
	"context"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

var (
	// contextTimeout is the global context timeout for requests to complete.
	contextTimeout = 15 * time.Second
)

// Provider returns the actual provider instance.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"googlecalendar_event": resourceEvent(),
		},
	}
	provider.ConfigureFunc = providerConfigure(provider)
	return provider
}

// providerConfigure configures the provider. Normally this would use schema
// data from the provider, but the provider loads all its configuration from the
// environment, so we just tell the config to load.
func providerConfigure(provider *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		var c Config
		if err := c.loadAndValidate(provider); err != nil {
			return nil, errors.Wrap(err, "failed to load config")
		}
		return &c, nil
	}
}

// contextWithTimeout creates a new context with the global context timeout.
func contextWithTimeout() (context.Context, func()) {
	return context.WithTimeout(context.Background(), contextTimeout)
}
