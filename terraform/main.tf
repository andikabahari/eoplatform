terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.41.0"
    }
  }

  backend "gcs" {
    bucket = "testenv-357307-ember"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = var.project
  region  = var.region
  zone    = var.zone
}

resource "google_storage_bucket" "ember" {
  name                        = "testenv-357307-ember"
  location                    = var.region
  uniform_bucket_level_access = true
}

module "cloud_run" {
  source  = "GoogleCloudPlatform/cloud-run/google"
  version = "~> 0.2.0"

  service_name = "eoplatform"
  project_id   = var.project
  location     = var.region
  image        = "docker.io/andikabahari/eoplatform"

  members = ["allUsers"]

  template_annotations = {
    "autoscaling.knative.dev/minScale" : 0,
    "autoscaling.knative.dev/maxScale" : 2,
    "generated-by" : "terraform",
    "run.googleapis.com/client-name" : "terraform"
  }

  requests = {
    "cpu" : "100m",
    "memory" : "128Mi",
  }

  limits = {
    "cpu" : "500m",
    "memory" : "512Mi",
  }

  env_vars = var.env_vars
}
