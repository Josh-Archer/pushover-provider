---
page_title: "Pushover Provider"
description: |-
  Use the Pushover provider to send push notifications and manage delivery groups via the Pushover API.
---

# Pushover Provider

The **Pushover** provider integrates with the [Pushover](https://pushover.net) push notification service. It allows you to send notifications and manage delivery group membership directly from Terraform or OpenTofu.

## Example Usage

```terraform
terraform {
  required_providers {
    pushover = {
      source  = "Josh-Archer/pushover"
      version = "~> 1.0"
    }
  }
}

provider "pushover" {
  api_token = var.pushover_api_token
}
```

The `api_token` can also be supplied via the `PUSHOVER_API_TOKEN` environment variable, which is the recommended approach for CI/CD pipelines:

```bash
export PUSHOVER_API_TOKEN="your_application_token"
terraform apply
```

## Authentication

You will need a **Pushover application API token**. Create one by registering an application at [https://pushover.net/apps/build](https://pushover.net/apps/build).

## Schema

### Required (one of)

- `api_token` (String, Sensitive) — Pushover application API token. Can also be provided via the `PUSHOVER_API_TOKEN` environment variable.

## Resources

- [pushover_message](resources/message.md) — Send a push notification.
- [pushover_group_user](resources/group_user.md) — Add a user to a Pushover delivery group.

## Data Sources

- [pushover_sounds](data-sources/sounds.md) — List available notification sounds.
- [pushover_validate_user](data-sources/validate_user.md) — Validate a user or group key.
