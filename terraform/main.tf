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
