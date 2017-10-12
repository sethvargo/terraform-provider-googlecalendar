# Terraform Google Calendar Provider

This is a [Terraform][terraform] provider for managing meetings on Google
Calendar. It enables you to treat "calendars as code" the same way you already
treat infrastructure as code!


## Installation

1. Download the latest compiled binary from [GitHub releases][releases].

1. Unzip/untar the archive.

1. Move it into `$HOME/.terraform.d/plugins`:

  ```sh
  $ mkdir -p $HOME/.terraform.d/plugins
  $ mv terraform-provider-googlecloud $HOME/.terraform.d/plugins/terraform-provider-googlecloud
  ```

1. Create your Terraform configurations as normal, and run `terraform init`:

  ```sh
  $ terraform init
  ```

  This will find the plugin locally.


## Usage

1. You will need a valid Google Cloud account and permission to create a service
account. You can create the service account in any project, but make sure you
choose "server-to-server" communication, since this is not an OAuth application.

  1. Visit the [Google Cloud Credentials Console][gcloud-creds]
  1. Click "Create credentials"
  1. Choose "Service account key"
  1. Use "Compute Engine default service account" (or make your own)
  1. Choose "JSON" as the key type
  1. Click "Create"

  After a few seconds, your browser will download a credentials file in JSON.
  Save this file securely (treat it like a password).

  Note: if you are making your own service account, be sure to grant access to
  the "calendar" OAuth scope!

1. Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to your newly-downloaded credentials file:

  ```sh
  $ export GOOGLE_APPLICATION_CREDENTIALS=/path/to/my-creds.json
  ```

  The Terraform provider automatically reads this environment variable and uses
  the file at the given path for authentication.

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

[terraform]: https://www.terraform.io/
[releases]: https://github.com/sethvargo/terraform-provider-googlecalendar/releases
[gcloud-creds]: https://console.cloud.google.com/apis/credentials
[examples]: https://github.com/sethvargo/terraform-provider-googlecalendar/tree/master/examples
