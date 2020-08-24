package googlecalendar

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Config is the structure used to instantiate the Google Calendar provider.
type Config struct {
	calendar *calendar.Service
}

// loadAndValidate loads the application default credentials from the
// environment and creates a client for communicating with Google APIs.
func (c *Config) loadAndValidate(provider *schema.Provider) error {
	log.Printf("[INFO] authenticating with local client")

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, calendar.CalendarScope)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	client.Transport = logging.NewTransport("Google", client.Transport)

	// Use a custom user-agent string. This helps google with analytics and it's
	// just a nice thing to do.
	userAgent := fmt.Sprintf("(%s %s) Terraform/%s",
		runtime.GOOS, runtime.GOARCH, provider.TerraformVersion)

	// Create the calendar service.
	calendarSvc, err := calendar.New(client)
	if err != nil {
		return nil
	}
	calendarSvc.UserAgent = userAgent
	c.calendar = calendarSvc

	return nil
}
