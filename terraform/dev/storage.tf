locals {
  bucket_name = "personal-website-blog-posts"
}

resource "google_storage_bucket" "posts" {
  location                    = var.region
  name                        = local.bucket_name
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  bucket  = google_storage_bucket.posts.name
  name    = "test.txt"
  content = "testing terraform storage object"
}