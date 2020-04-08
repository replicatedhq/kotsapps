#!/bin/bash

## Script to help set up GCP storage for snapshot functionality
## Requires: gsutil and for your account project to be linked via command line. 
## Validated on OSX shell

PROJECT_ID=$(gcloud config get-value project)
BUCKET_NAME=replsnap-$(LC_CTYPE=C tr -dc a-z0-9 < /dev/urandom | head -c 16 | xargs)
## Create the bucket
gsutil mb -l US gs://$BUCKET_NAME
## Create the iam account
gcloud iam service-accounts create "$BUCKET_NAME" --display-name="Velero service account"
## Set the $SERVICE_ACCOUNT_EMAIL variable to match `email` value
SERVICE_ACCOUNT_EMAIL=$(gcloud iam service-accounts list \
  --filter="displayName:Velero service account" \
  --format 'value(email)')

## Attach policies to give velero necessary permissions to function
ROLE_PERMISSIONS=(
    compute.disks.get
    compute.disks.create
    compute.disks.createSnapshot
    compute.snapshots.get
    compute.snapshots.create
    compute.snapshots.useReadOnly
    compute.snapshots.delete
    compute.zones.get
)

gcloud iam roles create velero.server \
    --project $PROJECT_ID \
    --title "Velero Server" \
    --permissions "$(IFS=","; echo "${ROLE_PERMISSIONS[*]}")"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member serviceAccount:$SERVICE_ACCOUNT_EMAIL \
    --role projects/$PROJECT_ID/roles/velero.server

gsutil iam ch serviceAccount:$SERVICE_ACCOUNT_EMAIL:objectAdmin gs://${BUCKET_NAME}

## Create the key linked to the IAM account and save
gcloud iam service-accounts keys create credentials-velero --iam-account=$SERVICE_ACCOUNT_EMAIL

## Copy key to Replicated Admin Console: /app/airgapped/snapshots/settings -- hycKNpgmi
pbcopy < "./credentials-velero" 
printf "\n\nBUCKET_NAME: $BUCKET_NAME\n\n"
echo "Enter bucket name, path, and paste key to Replicated Admin Console /app/airgapped/snapshots/settings"
