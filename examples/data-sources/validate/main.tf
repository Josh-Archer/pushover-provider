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

# Validate a user key and list their registered devices.
data "pushover_validate_user" "recipient" {
  user_key = var.pushover_user_key
}

output "is_group" {
  description = "Whether the key belongs to a Pushover delivery group."
  value       = data.pushover_validate_user.recipient.is_group
}

output "devices" {
  description = "Devices registered to this Pushover user."
  value       = data.pushover_validate_user.recipient.devices
}

output "licenses" {
  description = "Active Pushover licenses for this user."
  value       = data.pushover_validate_user.recipient.licenses
}

# Send to the user's first registered device.
resource "pushover_message" "first_device" {
  user_key = var.pushover_user_key
  message  = "Delivered to your first device: ${data.pushover_validate_user.recipient.devices[0]}"
  device   = data.pushover_validate_user.recipient.devices[0]
}
