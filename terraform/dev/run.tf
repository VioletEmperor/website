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


    containers {
      image = "${var.region}-docker.pkg.dev/${var.project}/${google_artifact_registry_repository.artifact_registry.name}/${local.cloud_run_service_image}"

      ports {
        container_port = 80
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

      env {
        name  = "PROJECT_ID"
        value = var.project
      }

      env {
        name = "FIREBASE_WEB_API_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.firebase_web_api_key.secret_id
            version = "latest"
          }
        }
      }

      env {
        name  = "STORAGE_MODE"
        value = "gcs"
      }

      env {
        name  = "GCS_BUCKET_NAME"
        value = google_storage_bucket.posts.name
      }

      env {
        name  = "GCS_PREFIX"
        value = "posts/"
      }

      env {
        name = "TURNSTILE_SECRET"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.turnstile_secret.secret_id
            version = "latest"
          }
        }
      }

    }
    service_account = google_service_account.service_account.email
  }

  depends_on = [
    google_secret_manager_secret_version.email_key_version,
    google_secret_manager_secret_version.firebase_web_api_key_version,
    google_secret_manager_secret_version.turnstile_secret_version,
    google_artifact_registry_repository.artifact_registry,
    google_firestore_database.database
  ]
}