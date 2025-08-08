variable "project" {
  description = "The id of the project where this module will be deployed"
  type        = string
}

variable "region" {
  description = "The region to host the resource in"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "The zone that the machine should be created in"
  type        = string
  default     = "us-central1-c"
}


variable "email_key" {
  description = "The email service API key"
  type        = string
  sensitive   = true
}

variable "firebase_web_api_key" {
  description = "Firebase Web API key for client-side authentication"
  type        = string
  sensitive   = true
}

