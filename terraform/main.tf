# vim:foldmethod=indent

provider "google" {
  version = "3.8.0"
  project = var.project_id
  region  = "us-east4"
  zone    = "us-east4-b"
}

variable "ssh_keys" {
  type        = string
  description = "The ssh username and key in the format `username:ssh-key`."
}

variable "project_id" {
  type        = string
  description = "The name of the google cloud project you're deploying to."
}

variable "deletion_protection" {
  type        = bool
  description = "Whether to apply deletion protection to the selfpass-instance."
  default     = true
}

output "selfpass_static_ip" {
  value       = google_compute_address.selfpass.address
  description = "The public static IP address used by the selfpass-instance."
}

data "google_compute_image" "debian_image" {
  family  = "debian-10"
  project = "debian-cloud"
}

resource "google_compute_address" "selfpass" {
  name = "selfpass-address"
}

resource "google_compute_disk" "selfpass_data" {
  name                      = "seflpass-data"
  type                      = "pd-ssd"
  snapshot                  = "selfpass-data-us-east4-c-20200215205733-qkksrkqr"
  physical_block_size_bytes = 4096
}

resource "google_compute_instance" "selfpass" {
  name                = "selfpass-instance"
  machine_type        = "f1-micro"
  deletion_protection = var.deletion_protection

  boot_disk {
    initialize_params {
      image = data.google_compute_image.debian_image.self_link
    }
  }

  attached_disk {
    source = google_compute_disk.selfpass_data.self_link
  }

  network_interface {
    network = "default"
    access_config {
      nat_ip = google_compute_address.selfpass.address
    }
  }

  metadata = {
    ssh-keys = var.ssh_keys
  }
}

resource "google_compute_firewall" "grpc_server" {
  name    = "grpc-server"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["8080"]
  }
}

resource "google_compute_firewall" "docker_machine" {
  name    = "docker-machine"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["2376"]
  }
}
