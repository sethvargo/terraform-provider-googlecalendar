// Package googlecalendar manages Google calendar events with Terraform.
package googlecalendar

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Provider returns the actual provider instance.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CREDENTIALS",
					"GOOGLE_CLOUD_KEYFILE_JSON",
					"GCLOUD_KEYFILE_JSON",
				}, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"googlecalendar_event": resourceEvent(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, provider, terraformVersion)
	}
	return provider
}

// providerConfigure configures the provider. Normally this would use schema
// data from the provider, but the provider loads all its configuration from the
// environment, so we just tell the config to load.
func providerConfigure(d *schema.ResourceData, p *schema.Provider, terraformVersion string) (interface{}, error) {
	var opts []option.ClientOption

	// Add credential source
	if v := d.Get("credentials").(string); v != "" {
		log.Printf("[TRACE] using supplied credentials")
		opts = append(opts, option.WithCredentialsJSON([]byte(v)))
	} else {
		log.Printf("[TRACE] attempting to use default credentials: %#v", d.Get("credentials"))
	}

	// Use a custom user-agent string. This helps google with analytics and it's
	// just a nice thing to do.
	userAgent := fmt.Sprintf("(%s %s) Terraform/%s",
		runtime.GOOS, runtime.GOARCH, terraformVersion)
	opts = append(opts, option.WithUserAgent(userAgent))

	log.Printf("[TRACE] client options: %v", opts)

	// Create the calendar service.
	ctx := context.Background()
	calendarSvc, err := calendar.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}
	calendarSvc.UserAgent = userAgent

	return &Config{
		calendar: calendarSvc,
	}, nil
}
