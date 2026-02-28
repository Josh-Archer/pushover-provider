# Configure the Pushover provider.
# The api_token can also be provided via the PUSHOVER_API_TOKEN environment variable.

terraform {
  required_providers {
    pushover = {
      source  = "Josh-Archer/pushover"
      version = "~> 1.0"
    }
  }
}

provider "pushover" {
  # Use a variable or environment variable – never hardcode credentials.
  api_token = var.pushover_api_token
}

variable "pushover_api_token" {
  description = "Pushover application API token."
  type        = string
  sensitive   = true
}
