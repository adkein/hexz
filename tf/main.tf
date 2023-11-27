terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "docker" {
  registry_auth {
    address = "europe-west6-docker.pkg.dev"
    config_file = pathexpand("~/.docker/config.json")
  }
}

provider "google" {
  credentials = file(var.credentials_file)

  project = var.project
  region  = var.region
  zone    = var.zone
}

resource "google_redis_instance" "cache" {
  name           = "hexz-game-state-cache"
  memory_size_gb = 1

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_artifact_registry_repository" "hexz-game-history-db" {
  location      = var.region
  repository_id = "hexz-game-history-db"
  description   = "hexz game history DB docker registry"
  format        = "DOCKER"
}

resource "docker_registry_image" "hexz-game-history-db" {
  name          = docker_image.hexz-game-history-db.name
}

resource "docker_image" "hexz-game-history-db" {
  name = "europe-west6-docker.pkg.dev/${var.project}/hexz-game-history-db/hexz-game-history-db"
  build {
    context = "${path.cwd}/../sql"
    tag     = ["europe-west6-docker.pkg.dev/${var.project}/hexz-game-history-db/hexz-game-history-db"]
  }
}

resource "google_service_account" "hexz-service-account" {
  account_id   = "hexz-provisioner"
  display_name = "hexz provisioner service account"
  description  = "initial provisioning of resources for hexz research and development"
}

resource "google_artifact_registry_repository_iam_member" "hexz-repo-iam" {
  location   = google_artifact_registry_repository.hexz-game-history-db.location
  repository = google_artifact_registry_repository.hexz-game-history-db.name
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:${google_service_account.hexz-service-account.email}"
}
