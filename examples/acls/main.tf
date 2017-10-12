// Create a google event and give guests (attendees) full control over the
// event
resource "googlecalendar_event" "demo" {
  summary     = "My Open Event"
  description = "Anyone can do anything, anytime, anywhere."
  location    = "Wherever you want!"

  start = "2017-10-12T15:00:00-05:00"
  end   = "2017-10-12T17:00:00-05:00"

  attendee {
    email = "seth@sethvargo.com"
  }

  // Allow guests to do anything
  guests_can_invite_others    = true
  guests_can_modify           = true
  guests_can_see_other_guests = true
}
