# Service Account Impersonation

This guide shows how to use the Terraform-created service account with gcloud CLI via impersonation.

## Setup

1. **Apply Terraform configuration** to grant impersonation permissions:
   ```bash
   cd terraform/prod
   terraform apply
   ```

2. **Set up impersonation**:
   ```bash
   gcloud config set auth/impersonate_service_account service-account@production-461100.iam.gserviceaccount.com
   ```

3. **Verify impersonation is active**:
   ```bash
   gcloud auth list
   # Should show the service account as impersonated
   ```

## Usage

Once impersonation is set up, all gcloud commands will use the service account:

```bash
# Deploy to Cloud Run
gcloud run deploy --source .

# Access Firestore
gcloud firestore databases list

# Push to Artifact Registry
gcloud builds submit --tag gcr.io/production-461100/myapp

# Access Cloud Storage
gsutil ls gs://your-bucket-name
```

## Stop Impersonation

To stop using the service account and return to your user account:

```bash
gcloud config unset auth/impersonate_service_account
```

## Service Account Permissions

The service account has these roles:
- `roles/datastore.user` - Firestore access
- `roles/artifactregistry.writer` - Push Docker images
- `roles/run.admin` - Manage Cloud Run services
- `roles/logging.logWriter` - Write application logs
- `roles/storage.objectUser` - Access storage objects
- `roles/secretmanager.secretAccessor` - Access secrets

## Troubleshooting

If you get permission errors:
1. Ensure Terraform was applied successfully
2. Check that your user account has the Token Creator role on the service account
3. Verify impersonation is active with `gcloud auth list`