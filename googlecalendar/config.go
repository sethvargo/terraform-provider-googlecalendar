package googlecalendar

import (
	"google.golang.org/api/calendar/v3"
)

// Config is the structure used to instantiate the Google Calendar provider.
type Config struct {
	calendar *calendar.Service
}
