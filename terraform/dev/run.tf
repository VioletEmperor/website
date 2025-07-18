locals {
  cloud_run_service_name  = "service"
  cloud_run_service_image = "cloudrun-image"
}

resource "google_cloud_run_v2_service" "default" {
  name                = "cloudrun-${local.cloud_run_service_name}"
  location            = var.region
  deletion_protection = false
  ingress             = "INGRESS_TRAFFIC_ALL"

  template {
    scaling {
      min_instance_count = 1
      max_instance_count = 1
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.instance.connection_name]
      }
    }

    containers {
      image = "${var.region}-docker.pkg.dev/${var.project}/${google_artifact_registry_repository.artifact_registry.name}/${local.cloud_run_service_image}"

      ports {
        container_port = 80
      }

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      env {
        name = "DB_HOST"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.database_host.secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "DB_USER"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.database_user.secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "DB_PASSWORD"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.database_password.secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "DB_NAME"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.database_name.secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "EMAIL_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.email_key.secret_id
            version = "latest"
          }
        }
      }
    }
    service_account = google_service_account.service_account.email
  }

  depends_on = [
    google_secret_manager_secret_version.database_host_version,
    google_secret_manager_secret_version.email_key_version
  ]
}