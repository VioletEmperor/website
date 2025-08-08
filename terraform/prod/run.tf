locals {
  cloud_run_service_name  = "prod-service"
  cloud_run_service_image = "cloudrun-image"
}

resource "google_cloud_run_v2_service" "default" {
  name                = "cloudrun-${local.cloud_run_service_name}"
  location            = var.region
  deletion_protection = false
  ingress             = "INGRESS_TRAFFIC_ALL"

  template {
    scaling {
      min_instance_count = 0
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

    }
    service_account = google_service_account.service_account.email
  }

  depends_on = [
    google_secret_manager_secret_version.email_key_version,
    google_secret_manager_secret_version.firebase_web_api_key_version,
    google_artifact_registry_repository.artifact_registry,
    google_firestore_database.database
  ]
}

# Domain mapping for the main domain (Cloudflare handles SSL)
# Commented out to avoid disrupting existing DNS during Terraform updates
# resource "google_cloud_run_domain_mapping" "domain" {
#   location = var.region
#   name     = "adamshkolnik.com"

#   metadata {
#     namespace = var.project
#   }

#   spec {
#     route_name       = google_cloud_run_v2_service.default.name
#     certificate_mode = "NONE"
#   }

#   depends_on = [google_cloud_run_v2_service.default]
# }

# Domain mapping for www subdomain (Cloudflare handles SSL)
# Commented out to avoid disrupting existing DNS during Terraform updates
# resource "google_cloud_run_domain_mapping" "www_domain" {
#   location = var.region
#   name     = "www.adamshkolnik.com"

#   metadata {
#     namespace = var.project
#   }

#   spec {
#     route_name       = google_cloud_run_v2_service.default.name
#     certificate_mode = "NONE"
#   }

#   depends_on = [google_cloud_run_v2_service.default]
# }