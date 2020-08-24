# Terraform Google Calendar Provider

This is a [Terraform][terraform] provider for managing meetings on Google
Calendar. It enables you to treat "calendars as code" the same way you already
treat infrastructure as code!


## Installation

1. Use Terraform 0.13+.

1. Pull from the module registry:

    ```hcl
    terraform {
      required_providers {
        googlecalendar = {
          source  = "sethvargo/googlecalendar"
          version = "~> 0.3"
        }
      }
    }
    ```

1. Create your Terraform configurations as normal, and run `terraform init`:

    ```sh
    $ terraform init
    ```


## Usage

1. Create a Terraform configuration file:

    ```hcl
    resource "googlecalendar_event" "example" {
      summary     = "My Event"
      description = "Long-form description of the event"
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
    }
    ```

1. Run `terraform init` to pull in the provider:

    ```sh
    $ terraform init
    ```

1. Run `terraform plan` and `terraform apply` to create events:

    ```sh
    $ terraform plan

    $ terraform apply
    ```

## Examples

For more examples, please see the [examples][examples] folder in this
repository.

## Reference

### Arguments

Arguments are provided as inputs to the resource, in the `*.tf` file.

- `summary` `(string, required)` - the "title" of the event.

- `start` `(string, required)` - the RFC3339-formatted start time of the event
  _with the timestamp included_.

- `end` `(string, required)` - the RFC3339-formatted end time of the event _with
  the timestamp included_.

- `description` `(string)` - the long-form description of the event. This can be
  multiple paragraphs using Terraform's heredoc syntax.

- `guests_can_invite_others` `(bool, true)` - specifies that guests (attendees)
  can invite other guests. Set this to false to allow only the organizer to
  manage the guest list.

- `guests_can_modify` `(bool, false)` - specifies that guests (attendees) can
  modify the event (change start time, description, etc). Set this to true to
  give any guest full control over the event.

- `guests_can_see_other_guests` `(bool, true)` - specifies that guests
  (attendees) can see other guests. Set this to false to restrict the guest list
  visibility.

- `show_as_available` `(bool, false)` - specifies that the time should be
  "blocked" on the calendar (mark as busy). Set this to true to create an event
  that is transparent.

- `send_notifications` `(bool, true)` - specifies that email notifications
  should be sent to guests (attendees). Set this to false to put things on
  people's calendar's without notifying them.

- `visibility` `(string)` - specifies the visibility for the event. Valid values
  are:

    - `""` - default inherit from calendar
    - `"public"` - public
    - `"private"` - private

- `attendee` `(list of structures)` - specifies a guest (attendee) to invite to
  the event. This may be specified more than once to invite multiple people to
  the same event. The following fields are supported:

    - `email` `(string, required)` - the Google email address of the attendee.

    - `optional` `(bool, false)` - specifies that the guest (attendee) is marked
      as optional. Set this to true to mark the user as an optional attendee.

- `reminder` `(list of reminders)` - specifies a reminder option leading up to
  the event for all attendees. This overrides any default reminders the user has
  set for their calendar. Leave this unset to inherit calendar default
  reminders. This may be specified more than once to remind multiple times. The
  following fields are supported:

    - `method` `(string, required)` - the method to use. Valid options are:

      - `"email"` - send an email
      - `"popup"` - popup in-browser
      - `"sms"` - send a text message (requires GSuite)

    - `before` `(string, required)` - the duration prior to the meeting to send
      the reminder. Note that the Google Calendar API expects this to be the
      "number of minutes", but Terraform adds syntactic sugar here by allowing
      you to specify the time in a Go timestamp like "30m" or "4h". These
      timestamps are parsed and converted to minutes automatically.

### Attributes

Attributes are values that are only known after creation.

- `event_id` `(string)` - the unique ID of the event on this calendar

- `hangout_link` `(string)` - the HTTPS web link to the attached Google Hangout.
  In practice, I have been unable to get this link to appear.

- `html_link` `(string)` - the HTTP web link to the calendar invite on
  [calendar.google.com](https://calendar.google.com/).


## Constraints & Understanding

### Architecture

It is important to note that you are not creating an event on a person's
calendar - that is not permitted. Can you imagine if anyone with a Google Cloud
account could put events on your calendar!? Instead, you are creating an event
on the service account's calendar and then inviting the appropriate attendees to
that event. In this way, Terraform acts as a "robot" which invites people to an
event.

### Time

A good future enhancement is to allow "human" times. Right now, all times _must_
be specified in RFC3339 format. It would be great to allow arbitrary human times
like "Oct 13, 2017 at 4pm EST".

## License & Author

This project is licensed under the MIT license by Seth Vargo
(seth@sethvargo.com).

[terraform]: https://www.terraform.io/
[releases]: https://github.com/sethvargo/terraform-provider-googlecalendar/releases
[gcloud-creds]: https://console.cloud.google.com/apis/credentials
[examples]: https://github.com/sethvargo/terraform-provider-googlecalendar/tree/master/examples
