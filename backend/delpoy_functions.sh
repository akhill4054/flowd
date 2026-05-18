#!/usr/bin/env bash

set -e

source .env
source .env.local

FUNCTION="${1:-all}"
ENV_VARS=$(grep -v '^#' .env | xargs | tr ' ' ',')

gcloud config set project "$GCP_PROJECT_ID"

deploy() {
  local name="$1"
  local entry="$2"

  echo "Deploying $name..."

  gcloud functions deploy "$name" \
    --project="$GCP_PROJECT_ID" \
    --gen2 \
    --runtime=go124 \
    --region=asia-south1 \
    --source=. \
    --entry-point="$entry" \
    --trigger-http \
    --set-env-vars="$ENV_VARS" \
    --allow-unauthenticated
}

case "$FUNCTION" in
  meta-webhook-handler)
    deploy flowd-meta-webhook-handler MetaWebhookHandler
    ;;

  dashboard-api)
    deploy flowd-dashboard-api DashboardApiFunction
    ;;

  all)
    deploy flowd-meta-webhook-handler MetaWebhookHandler
    deploy flowd-dashboard-api DashboardApiFunction
    ;;

  *)
    echo "Usage: ./deploy.sh [meta-webhook-handler|dashboard-api|all]"
    exit 1
    ;;
esac
