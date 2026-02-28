---
page_title: "pushover_sounds Data Source - pushover"
subcategory: ""
description: |-
  Retrieves the list of notification sounds available to your Pushover application.
---

# pushover_sounds (Data Source)

Retrieves all notification sounds that are available to your Pushover application. Use the returned keys as the `sound` argument of `pushover_message`.

## Example Usage

### List all sounds

```terraform
data "pushover_sounds" "all" {}

output "available_sounds" {
  value = data.pushover_sounds.all.sounds
}

output "sound_keys" {
  value = data.pushover_sounds.all.keys
}
```

### Use a sound in a message

```terraform
data "pushover_sounds" "all" {}

resource "pushover_message" "alert" {
  user_key = var.pushover_user_key
  message  = "Disk usage exceeded 90%"
  # Pick a sound that exists in the current application's sound list.
  sound    = contains(data.pushover_sounds.all.keys, "siren") ? "siren" : "pushover"
}
```

## Schema

### Read-Only

- `id` (String) — Placeholder identifier (`sounds`).
- `keys` (List of String) — All available sound keys. Pass one to `pushover_message.sound`.
- `sounds` (Map of String) — Map of sound key → human-readable name (e.g., `"pushover"` → `"Pushover (default)"`).
