locals {
  id = "docker-registry"
}

resource "google_artifact_registry_repository" "artifact_registry" {
  format        = "DOCKER"
  repository_id = local.id
  location      = var.region
  description   = "image registry for cloud run docker image"
}