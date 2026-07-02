# Shared across every scenario root module via a committed symlink
# (scenarios/<name>/_shared.tf -> ../_shared/_shared.tf).
#
# It holds the terraform{}/provider blocks and the SUPERSET of all input
# variable declarations. Declaring every variable in every scenario keeps a
# single shared secret.tfvars usable by all scenarios with zero "Value for
# undeclared variable" warnings (declared-but-unused variables are not an error).
# When a new scenario needs a new credential variable, declare it here once.
#
# Only the rudderstack provider is referenced, and it is dev-overridden by
# run.sh, so `terraform init` resolves the local module offline (no registry).

terraform {
  required_providers {
    rudderstack = { source = "rudderstack.com/rudderlabs/rudderstack" }
  }
}

provider "rudderstack" {
  api_url      = var.api_url
  access_token = var.access_token
}

variable "access_token" {
  description = "RudderStack staging personal access token."
  type        = string
  sensitive   = true
}

variable "api_url" {
  description = "RudderStack API base URL."
  type        = string
  default     = "https://api.staging.rudderlabs.com"
}

variable "bq_project" {
  description = "GCP project ID where the BigQuery dataset lives."
  type        = string
}

variable "bq_location" {
  description = "BigQuery dataset location (e.g. US, EU)."
  type        = string
  default     = "US"
}

variable "bq_dataset" {
  description = "BigQuery dataset (schema) name."
  type        = string
}

variable "bq_table" {
  description = "BigQuery table name."
  type        = string
}

variable "bq_credentials" {
  description = "GCP service-account key JSON (contents of the JSON key file)."
  type        = string
  sensitive   = true
}

# ── Customer.io (VDM v2) scenario ────────────────────────────────────────────
variable "customerio_site_id" {
  description = "Customer.io site ID."
  type        = string
  default     = ""
  sensitive   = true
}

variable "customerio_api_key" {
  description = "Customer.io API key."
  type        = string
  default     = ""
  sensitive   = true
}

variable "customerio_datacenter" {
  description = "Customer.io data center (US or EU)."
  type        = string
  default     = "US"
}

# ── Customer.io Audience scenario ────────────────────────────────────────────
variable "customerio_audience_app_api_key" {
  description = "Customer.io App API key for the Audience destination."
  type        = string
  default     = ""
  sensitive   = true
}

variable "customerio_audience_region" {
  description = "Customer.io Audience destination region (US or EU)."
  type        = string
  default     = "US"
}

variable "customerio_audience_id" {
  description = "Customer.io audience ID (positive integer)."
  type        = number
  default     = 0
}
