
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

# TURNSTILE_SECRET
resource "google_secret_manager_secret" "turnstile_secret" {
  secret_id = "turnstile-secret"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "turnstile_secret_version" {
  secret      = google_secret_manager_secret.turnstile_secret.id
  secret_data = var.turnstile_secret
}

resource "google_secret_manager_secret_iam_member" "turnstile_secret_policy" {
  secret_id  = google_secret_manager_secret.turnstile_secret.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.service_account.email}"
  depends_on = [google_secret_manager_secret.turnstile_secret]
}

