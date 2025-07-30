locals {
  database_version = "POSTGRES_17"
}

resource "google_sql_database_instance" "instance" {
  name                = "prod-cloudsql-instance-${var.database_name}"
  region              = var.region
  database_version    = local.database_version
  deletion_protection = true

  settings {
    tier                = "db-f1-micro"
    availability_type   = "ZONAL" 
    disk_type          = "PD_HDD"
    disk_size          = 10
    disk_autoresize    = false
    
    backup_configuration {
      enabled = false
    }

    ip_configuration {
      ipv4_enabled = true
      require_ssl  = false
    }
  }
}

resource "google_sql_database" "instance" {
  instance = google_sql_database_instance.instance.name
  name     = var.database_name
}

resource "google_sql_user" "user" {
  instance = google_sql_database_instance.instance.name
  name     = var.database_user
  password = var.database_password
}