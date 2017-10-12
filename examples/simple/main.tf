// Create a google calendar event.
resource "googlecalendar_event" "demo" {
  // Common options
  summary     = "My Demo Terraform Event"
  description = "Long-form description of the event, such as why it's needed"
  location    = "Conference Room B"

  // Start and end times work best if specified as RFC3339.
  start = "2017-10-12T15:00:00-05:00"
  end   = "2017-10-12T17:00:00-05:00"

  // Each attendee is listed separately, and attendees can be marked as
  // optional.
  attendee {
    email = "seth@sethvargo.com"
  }

  attendee {
    email    = "you@company.com"
    optional = true
  }

  // By default, the user's calendar reminders are used. By setting any
  // reminders, you override all default calendar reminders. The Google API
  // expects calendar  reminder times to be in "minutes", but you can specify
  // them as a Go-style time.Duration value for simplicity here, like "30m" for
  // "30 minutes".
  reminder {
    method = "email"
    before = "1h"
  }

  reminder {
    method = "popup"
    before = "5m"
  }
}
