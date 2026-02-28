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

variable "pushover_api_token" {
  type      = string
  sensitive = true
}

variable "pushover_user_key" {
  type      = string
  sensitive = true
}

# Fetch all available sounds.
data "pushover_sounds" "all" {}

# Display all sound keys.
output "sound_keys" {
  description = "All sound keys available to this application."
  value       = data.pushover_sounds.all.keys
}

# Display the full sound map (key → name).
output "sounds_map" {
  description = "Map of sound key to human-readable name."
  value       = data.pushover_sounds.all.sounds
}

# Use a dynamically fetched sound in a message.
resource "pushover_message" "magic_sound" {
  user_key = var.pushover_user_key
  message  = "Notification with the 'magic' sound."
  sound    = contains(data.pushover_sounds.all.keys, "magic") ? "magic" : "pushover"
}
