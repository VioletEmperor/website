locals {
  database_version = "POSTGRES_17"
}

resource "google_sql_database_instance" "instance" {
  name                = "cloudsql-instance-${var.database_name}"
  region              = var.region
  database_version    = local.database_version
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_database" "instance" {
  instance = google_sql_database_instance.instance.name
  name     = "cloudsql-${var.database_name}"
}

resource "google_sql_user" "user" {
  instance = google_sql_database_instance.instance.name
  name     = var.database_user
  password = var.database_password
}

resource "google_sql_da" "" {}