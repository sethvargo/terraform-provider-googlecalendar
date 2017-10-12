package googlecalendar

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

var oauthScopes = []string{
	calendar.CalendarScope,
}

// Config is the structure used to instantiate the Google Calendar provider.
type Config struct {
	calendar *calendar.Service
}

// loadAndValidate loads the application default credentials from the
// environment and creates a client for communicating with Google APIs.
func (c *Config) loadAndValidate() error {
	log.Printf("[INFO] authenticating with local client")
	client, err := google.DefaultClient(context.Background(), oauthScopes...)
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	// Use a custom user-agent string. This helps google with analytics and it's
	// just a nice thing to do.
	client.Transport = logging.NewTransport("Google", client.Transport)
	userAgent := fmt.Sprintf("(%s %s) Terraform/%s",
		runtime.GOOS, runtime.GOARCH, terraform.VersionString())

	// Create the calendar service.
	calendarSvc, err := calendar.New(client)
	if err != nil {
		return nil
	}
	calendarSvc.UserAgent = userAgent
	c.calendar = calendarSvc

	return nil
}
