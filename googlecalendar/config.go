package googlecalendar

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/hashicorp/terraform/helper/schema"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Config is the structure used to instantiate the Google Calendar provider.
type Config struct {
	calendar *calendar.Service
}

// loadAndValidate loads the application default credentials from the
// environment and creates a client for communicating with Google APIs.
func (c *Config) loadAndValidate(provider *schema.Provider) error {
	log.Printf("[INFO] authenticating with local client")

	// Use a custom user-agent string. This helps google with analytics and it's
	// just a nice thing to do.
	userAgent := fmt.Sprintf("(%s %s) Terraform/%s",
		runtime.GOOS, runtime.GOARCH, provider.TerraformVersion)

	// Create the calendar service.
	ctx := context.Background()
	calendarSvc, err := calendar.NewService(ctx,
		option.WithUserAgent(userAgent))
	if err != nil {
		return nil
	}
	c.calendar = calendarSvc

	return nil
}
