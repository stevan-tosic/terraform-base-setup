variable "project" {}
variable "credentials_file" {}

variable "region" {
  type        = string
  default     = "europe-west6"
  description = "Regions are independent geographic areas that consist of zones."
}

variable "zone" {
  type        = string
  default     = "europe-west6-c"
  description = "Deployment area for Google Cloud resources within a region"
}

variable "os_image" {
  type        = string
  default     = "debian-cloud/debian-11"
  description = "Operating system image"
}

variable "vm_params" {
  type = object({
    name         = string,
    machine_type = string,
    zone         = string,

    allow_stopping_for_update = bool,
  })
  default = {
    name         = "terraform-instance",
    machine_type = "f1-micro",
    zone         = "europe-west6-b",

    allow_stopping_for_update = true,
  }
  description = "vm parameters"

  validation {
    condition     = length(var.vm_params.name) > 3
    error_message = "VM name must be at least 4 characters."
  }
}