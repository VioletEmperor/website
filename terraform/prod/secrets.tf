# DB_HOST
resource "google_secret_manager_secret" "database_host" {
  secret_id = "database-host"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "database_host_version" {
  secret      = google_secret_manager_secret.database_host.id
  secret_data = "/cloudsql/${var.project}:${var.region}:${google_sql_database_instance.instance.id}"
}

resource "google_secret_manager_secret_iam_member" "database_host_policy" {
  secret_id  = google_secret_manager_secret.database_host.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.service_account.email}"
  depends_on = [google_secret_manager_secret.database_host]
}

# DB_USER
resource "google_secret_manager_secret" "database_user" {
  secret_id = "database-user"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "database_user_version" {
  secret      = google_secret_manager_secret.database_user.id
  secret_data = var.database_user
}

resource "google_secret_manager_secret_iam_member" "database_user_policy" {
  secret_id  = google_secret_manager_secret.database_user.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.service_account.email}"
  depends_on = [google_secret_manager_secret.database_user]
}

# DB_PASSWORD
resource "google_secret_manager_secret" "database_password" {
  secret_id = "database-password"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "database_password_version" {
  secret      = google_secret_manager_secret.database_password.id
  secret_data = var.database_password
}

resource "google_secret_manager_secret_iam_member" "database_password_policy" {
  secret_id  = google_secret_manager_secret.database_password.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.service_account.email}"
  depends_on = [google_secret_manager_secret.database_password]
}

# DB_NAME
resource "google_secret_manager_secret" "database_name" {
  secret_id = "database-name"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "database_name_version" {
  secret      = google_secret_manager_secret.database_name.id
  secret_data = var.database_name
}

resource "google_secret_manager_secret_iam_member" "database_name_policy" {
  secret_id  = google_secret_manager_secret.database_name.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.service_account.email}"
  depends_on = [google_secret_manager_secret.database_name]
}

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

