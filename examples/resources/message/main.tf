# --- Variables ---
variable "pushover_api_token" {
  description = "Pushover application API token."
  type        = string
  sensitive   = true
}

variable "pushover_user_key" {
  description = "Pushover user or group key to send notifications to."
  type        = string
  sensitive   = true
}

# --- Provider ---
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

# --- Example 1: Simple message ---
resource "pushover_message" "simple" {
  user_key = var.pushover_user_key
  message  = "Hello from Terraform!"
}

# --- Example 2: Rich message with all common fields ---
resource "pushover_message" "rich" {
  user_key  = var.pushover_user_key
  message   = "Deploy of <b>my-app v2.3.0</b> completed successfully in 42 s."
  title     = "✅ Deploy Complete"
  url       = "https://my-app.example.com"
  url_title = "Open Application"
  priority  = 1
  sound     = "magic"
  html      = true
  ttl       = 86400  # Auto-delete from Pushover servers after 24 h
}

# --- Example 3: Emergency message ---
resource "pushover_message" "emergency" {
  user_key = var.pushover_user_key
  message  = "Production database is DOWN!"
  title    = "🔴 CRITICAL OUTAGE"
  priority = 2
  retry    = 60     # Resend every minute
  expire   = 3600   # Give up after 1 hour
  callback = "https://ops.example.com/webhook/ack"
}

# Output the emergency receipt so it can be polled later.
output "emergency_receipt" {
  description = "Pushover receipt token for the emergency notification."
  value       = pushover_message.emergency.receipt
  sensitive   = false
}

# --- Example 4: Low-priority quiet notification ---
resource "pushover_message" "quiet" {
  user_key = var.pushover_user_key
  message  = "Nightly backup completed at 03:00 UTC."
  title    = "Backup Report"
  priority = -1
  sound    = "none"
}

# --- Example 5: Targeted device ---
resource "pushover_message" "device_specific" {
  user_key = var.pushover_user_key
  message  = "This notification goes only to your iPhone."
  device   = "iphone"
}
