provider "google" {
  project     = var.project
  credentials = file(var.credentials_file)
  region      = var.region
  zone        = var.zone
}

resource "google_compute_instance" "st_instance" {
  machine_type = var.vm_params.machine_type
  name         = var.vm_params.name
  zone         = var.vm_params.zone

  allow_stopping_for_update = var.vm_params.allow_stopping_for_update

  boot_disk {
    initialize_params {
      image = var.os_image
    }
  }
  network_interface {
    network    = google_compute_network.terraform_network.self_link
    subnetwork = google_compute_subnetwork.terraform_subnet.self_link
    access_config {
      // necessary even empty
    }
  }
}

resource "google_compute_network" "terraform_network" {
  name                    = "terraform-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "terraform_subnet" {
  ip_cidr_range = "10.20.0.0/16"
  name          = "terraform-subnetwork"
  region        = "europe-west6"
  network       = google_compute_network.terraform_network.id
}