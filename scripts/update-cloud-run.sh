#!/bin/bash

region=us-central1
service=cloudrun-service
project=develoment-461100
repository=docker-registry
image=cloudrun-image

gcloud builds submit --region=${region} \
  --tag ${region}-docker.pkg.dev/${project}/${repository}/${image}

gcloud run deploy ${service} \
  --image ${region}-docker.pkg.dev/${project}/${repository}/${image}

gcloud run services update ${service} --no-invoker-iam-check