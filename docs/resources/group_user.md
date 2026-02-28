---
page_title: "pushover_group_user Resource - pushover"
subcategory: ""
description: |-
  Adds a user to a Pushover delivery group and manages their membership settings.
---

# pushover_group_user (Resource)

Adds a Pushover user to a delivery group. The group must already exist (create it in the [Pushover dashboard](https://pushover.net) or via your Pushover subscription API). This resource manages a single group–user relationship.

Changing `group_key`, `user_key`, or `device` forces a new resource. Changing `memo` or `disabled` is updated in-place.

## Example Usage

### Basic membership

```terraform
resource "pushover_group_user" "ops_on_call" {
  group_key = var.pushover_group_key
  user_key  = var.pushover_user_key
  memo      = "Primary on-call engineer"
}
```

### Device-scoped membership

```terraform
resource "pushover_group_user" "mobile_only" {
  group_key = var.pushover_group_key
  user_key  = var.pushover_user_key
  device    = "iphone"
  memo      = "Receives on-call alerts on iPhone only"
}
```

### Managing multiple members

```terraform
variable "on_call_users" {
  type = map(object({
    user_key = string
    memo     = string
  }))
  default = {
    alice = { user_key = "uAliceKey", memo = "Primary"   }
    bob   = { user_key = "uBobKey",   memo = "Secondary" }
  }
}

resource "pushover_group_user" "team" {
  for_each  = var.on_call_users
  group_key = var.pushover_group_key
  user_key  = each.value.user_key
  memo      = each.value.memo
}
```

### Temporarily disabling a member

```terraform
resource "pushover_group_user" "engineer" {
  group_key = var.pushover_group_key
  user_key  = var.pushover_user_key
  disabled  = var.engineer_on_leave   # set to true during leave
}
```

## Schema

### Required

- `group_key` (String) — The Pushover delivery group key. **(Forces replacement)**
- `user_key` (String) — The Pushover user key to add to the group. **(Forces replacement)**

### Optional

- `device` (String) — Restrict notifications to this specific device for the user. **(Forces replacement)**
- `disabled` (Boolean) — Set to `true` to disable notifications without removing the user from the group. Default: `false`.
- `memo` (String) — A note about this group member (≤ 200 characters).

### Read-Only

- `id` (String) — Computed unique identifier: `group_key/user_key` or `group_key/user_key/device` when a device is specified.

## Import

Group user resources can be imported using the format `group_key/user_key` or `group_key/user_key/device`:

```shell
terraform import pushover_group_user.ops_on_call gYourGroupKey/uYourUserKey
terraform import pushover_group_user.mobile_only gYourGroupKey/uYourUserKey/iphone
```
