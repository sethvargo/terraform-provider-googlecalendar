package googlecalendar

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pkg/errors"
	calendar "google.golang.org/api/calendar/v3"
)

var (
	eventValidMethods = []string{"email", "popup", "sms"}

	eventValidVisbilities = []string{"public", "private"}
)

func resourceEvent() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventCreate,
		Read:   resourceEventRead,
		Update: resourceEventUpdate,
		Delete: resourceEventDelete,

		Schema: map[string]*schema.Schema{
			"summary": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"location": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"start": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"end": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"guests_can_invite_others": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"guests_can_modify": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"guests_can_see_other_guests": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"show_as_available": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"send_notifications": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"visibility": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.StringInSlice(eventValidVisbilities, false),
			},

			"attendee": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"optional": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"reminder": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(eventValidMethods, false),
						},

						"before": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			//
			// Computed values
			//
			"event_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"hangout_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"html_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceEventCreate creates a new event via the API.
func resourceEventCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	event, err := resourceEventBuild(d, meta)
	if err != nil {
		return errors.Wrap(err, "failed to build event")
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()
	eventAPI, err := config.calendar.Events.
		Insert("primary", event).
		SendNotifications(d.Get("send_notifications").(bool)).
		MaxAttendees(25).
		Context(ctx).
		Do()
	if err != nil {
		return errors.Wrap(err, "failed to create event")
	}

	d.SetId(eventAPI.Id)

	return resourceEventRead(d, meta)
}

// resourceEventRead reads information about the event from the API.
func resourceEventRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := contextWithTimeout()
	defer cancel()
	event, err := config.calendar.Events.
		Get("primary", d.Id()).
		Context(ctx).
		Do()
	if err != nil {
		return errors.Wrap(err, "failed to read event")
	}

	d.Set("summary", event.Summary)
	d.Set("location", event.Location)
	d.Set("description", event.Description)
	d.Set("start", event.Start)
	d.Set("end", event.End)
	if event.GuestsCanInviteOthers != nil {
		d.Set("guests_can_invite_others", *event.GuestsCanInviteOthers)
	}
	d.Set("guests_can_modify", event.GuestsCanModify)
	if event.GuestsCanSeeOtherGuests != nil {
		d.Set("guests_can_see_other_guests", *event.GuestsCanSeeOtherGuests)
	}
	d.Set("show_as_available", transparencyToBool(event.Transparency))
	d.Set("visibility", event.Visibility)

	// Handle reminders
	if event.Reminders != nil && len(event.Reminders.Overrides) > 0 {
		d.Set("reminder", flattenEventReminders(event.Reminders.Overrides))
	}

	// Handle attendees
	if len(event.Attendees) > 0 {
		d.Set("attendee", flattenEventAttendees(event.Attendees))
	}

	d.Set("event_id", event.Id)
	d.Set("hangout_link", event.HangoutLink)
	d.Set("html_link", event.HtmlLink)

	return nil
}

// resourceEventUpdate updates an event via the API.
func resourceEventUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	event, err := resourceEventBuild(d, meta)
	if err != nil {
		return errors.Wrap(err, "failed to build event")
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()
	eventAPI, err := config.calendar.Events.
		Update("primary", d.Id(), event).
		SendNotifications(d.Get("send_notifications").(bool)).
		MaxAttendees(25).
		Context(ctx).
		Do()
	if err != nil {
		return errors.Wrap(err, "failed to update event")
	}

	d.SetId(eventAPI.Id)

	return resourceEventRead(d, meta)
}

// resourceEventDelete deletes an event via the API.
func resourceEventDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id := d.Id()
	sendNotifications := d.Get("send_notifications").(bool)

	ctx, cancel := contextWithTimeout()
	defer cancel()
	err := config.calendar.Events.
		Delete("primary", id).
		SendNotifications(sendNotifications).
		Context(ctx).
		Do()
	if err != nil {
		return errors.Wrap(err, "failed to delete event")
	}

	d.SetId("")

	return nil
}

// resourceBuildEvent is a shared helper function which builds an "event" struct
// from the schema. This is used by create and update.
func resourceEventBuild(d *schema.ResourceData, meta interface{}) (*calendar.Event, error) {
	summary := d.Get("summary").(string)
	location := d.Get("location").(string)
	description := d.Get("description").(string)

	start := d.Get("start").(string)
	end := d.Get("end").(string)

	guestsCanInviteOthers := d.Get("guests_can_invite_others").(bool)
	guestsCanModify := d.Get("guests_can_modify").(bool)
	guestsCanSeeOtherGuests := d.Get("guests_can_see_other_guests").(bool)
	showAsAvailable := d.Get("show_as_available").(bool)
	visibility := d.Get("visibility").(string)

	var event calendar.Event
	event.Summary = summary
	event.Location = location
	event.Description = description
	event.GuestsCanInviteOthers = &guestsCanInviteOthers
	event.GuestsCanModify = guestsCanModify
	event.GuestsCanSeeOtherGuests = &guestsCanSeeOtherGuests
	event.Source = &calendar.EventSource{
		Title: "Terraform by HashiCorp",
		Url:   "https://www.terraform.io/",
	}
	event.Transparency = boolToTransparency(showAsAvailable)
	event.Visibility = visibility
	event.Start = &calendar.EventDateTime{
		DateTime: start,
	}
	event.End = &calendar.EventDateTime{
		DateTime: end,
	}

	// Parse reminders
	remindersRaw := d.Get("reminder").(*schema.Set)
	if remindersRaw.Len() > 0 {
		reminders := make([]*calendar.EventReminder, remindersRaw.Len())

		for i, v := range remindersRaw.List() {
			m := v.(map[string]interface{})

			d, err := time.ParseDuration(m["before"].(string))
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse 'before'")
			}
			minutes := int64(d.Round(time.Minute).Minutes())

			reminders[i] = &calendar.EventReminder{
				Method:  m["method"].(string),
				Minutes: minutes,
			}
		}

		event.Reminders = &calendar.EventReminders{
			Overrides:       reminders,
			UseDefault:      false,
			ForceSendFields: []string{"UseDefault"},
		}
	}

	// Parse attendees
	attendeesRaw := d.Get("attendee").(*schema.Set)
	if attendeesRaw.Len() > 0 {
		attendees := make([]*calendar.EventAttendee, attendeesRaw.Len())

		for i, v := range attendeesRaw.List() {
			m := v.(map[string]interface{})

			attendees[i] = &calendar.EventAttendee{
				Email:    m["email"].(string),
				Optional: m["optional"].(bool),
			}
		}

		event.Attendees = attendees
	}

	return &event, nil
}

// flattenEventAttendees flattens the list of event reminders into a map for
// storing in the schema.
func flattenEventAttendees(list []*calendar.EventAttendee) []map[string]interface{} {
	result := make([]map[string]interface{}, len(list))
	for i, v := range list {
		result[i] = map[string]interface{}{
			"email":    v.Email,
			"optional": v.Optional,
		}
	}
	return result
}

// flattenEventReminders flattens the list of event reminders into a map for
// storing in the schema.
func flattenEventReminders(list []*calendar.EventReminder) []map[string]interface{} {
	result := make([]map[string]interface{}, len(list))
	for i, v := range list {
		result[i] = map[string]interface{}{
			"method": v.Method,
			"before": fmt.Sprintf("%dm", v.Minutes),
		}
	}
	return result
}

// boolToTransparency converts a boolean representing "show as available" to the
// corresponding transpency string.
func boolToTransparency(showAsAvailable bool) string {
	if !showAsAvailable {
		return "opaque"
	}
	return "transparent"
}

// transparencyToBool converts a transparency string into a boolean representing
// "show as available".
func transparencyToBool(s string) bool {
	switch s {
	case "opaque":
		return false
	case "transparent":
		return true
	default:
		log.Printf("[WARN] unknown transparency %q", s)
		return false
	}
}
