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
