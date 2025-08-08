
# EMAIL_KEY
resource "google_secret_manager_secret" "email_key" {
  secret_id = "email-key"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "email_key_version" {
  secret      = google_secret_manager_secret.email_key.id
  secret_data = var.email_key
}

resource "google_secret_manager_secret_iam_member" "email_key_policy" {
  secret_id  = google_secret_manager_secret.email_key.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.service_account.email}"
  depends_on = [google_secret_manager_secret.email_key]
}

# FIREBASE_WEB_API_KEY
resource "google_secret_manager_secret" "firebase_web_api_key" {
  secret_id = "firebase-web-api-key"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "firebase_web_api_key_version" {
  secret      = google_secret_manager_secret.firebase_web_api_key.id
  secret_data = var.firebase_web_api_key
}

resource "google_secret_manager_secret_iam_member" "firebase_web_api_key_policy" {
  secret_id  = google_secret_manager_secret.firebase_web_api_key.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.service_account.email}"
  depends_on = [google_secret_manager_secret.firebase_web_api_key]
}

