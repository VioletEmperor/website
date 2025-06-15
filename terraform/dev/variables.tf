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

variable "database_name" {
  description = "The name of the database to connect to"
  type        = string
  sensitive   = true
}

variable "database_user" {
  description = "The name of the default database user that will be created on the database"
  type        = string
  sensitive   = true
}

variable "database_password" {
  description = "The password of the default database user that will be created on the database"
  type        = string
  sensitive   = true
}