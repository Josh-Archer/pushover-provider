---
page_title: "pushover_validate_user Data Source - pushover"
subcategory: ""
description: |-
  Validates a Pushover user or group key and returns the associated devices and licenses.
---

# pushover_validate_user (Data Source)

Validates a Pushover user or group key against the Pushover API. Returns the registered device names and active licenses for the key.

This data source is useful for:
- Verifying that a recipient key is valid before sending notifications.
- Enumerating available device names to use with `pushover_message.device`.
- Detecting whether a key belongs to a delivery group (`is_group = true`).

## Example Usage

### Validate a user key

```terraform
data "pushover_validate_user" "recipient" {
  user_key = var.pushover_user_key
}

output "recipient_is_group" {
  value = data.pushover_validate_user.recipient.is_group
}

output "recipient_devices" {
  value = data.pushover_validate_user.recipient.devices
}
```

### Validate and target a specific device

```terraform
data "pushover_validate_user" "user" {
  user_key = var.pushover_user_key
}

resource "pushover_message" "first_device" {
  user_key = var.pushover_user_key
  message  = "Delivered to your first registered device."
  # Use the first device from the validated list.
  device   = data.pushover_validate_user.user.devices[0]
}
```

### Validate with a specific application token

```terraform
data "pushover_validate_user" "check" {
  user_key  = var.pushover_user_key
  api_token = var.app_specific_token
}
```

## Schema

### Required

- `user_key` (String) — The Pushover user or group key to validate.

### Optional

- `api_token` (String, Sensitive) — Override the provider-level API token for this validation.
- `device` (String) — Restrict validation to a specific device name.

### Read-Only

- `id` (String) — The validated user key (same as `user_key`).
- `devices` (List of String) — Device names registered to this user.
- `is_group` (Boolean) — `true` if this key belongs to a Pushover delivery group.
- `licenses` (List of String) — Active license types (e.g., `"iOS"`, `"Android"`).
