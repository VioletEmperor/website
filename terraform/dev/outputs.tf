# Cloud Run Service
output "cloud_run_url" {
  description = "URL of the deployed Cloud Run service"
  value       = google_cloud_run_v2_service.default.uri
}

output "cloud_run_service_name" {
  description = "Name of the Cloud Run service"
  value       = google_cloud_run_v2_service.default.name
}

# Firestore Database Information
output "firestore_database_name" {
  description = "Firestore database name"
  value       = google_firestore_database.database.name
}

# Storage Information
output "storage_bucket_name" {
  description = "Name of the GCS bucket for blog posts"
  value       = google_storage_bucket.posts.name
}

output "storage_bucket_url" {
  description = "GCS bucket URL"
  value       = google_storage_bucket.posts.url
}

# Service Account
output "service_account_email" {
  description = "Email of the service account used by Cloud Run"
  value       = google_service_account.service_account.email
}

# Artifact Registry
output "artifact_registry_repository" {
  description = "Artifact Registry repository name"
  value       = google_artifact_registry_repository.artifact_registry.name
}

output "docker_image_url" {
  description = "Full Docker image URL for deployment"
  value       = "${var.region}-docker.pkg.dev/${var.project}/${google_artifact_registry_repository.artifact_registry.name}/${local.cloud_run_service_image}"
}