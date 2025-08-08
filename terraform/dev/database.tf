# Firestore database
resource "google_firestore_database" "database" {
  project         = var.project
  name            = "(default)"
  location_id     = "us-central"
  type            = "FIRESTORE_NATIVE"
  deletion_policy = "DELETE"
}