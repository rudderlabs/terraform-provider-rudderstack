# Inputs for the shared BigQuery rETL source module. Each scenario instantiates
# this module with a unique name_prefix so the created account/source names don't
# collide across scenarios.
variable "name_prefix" {
  description = "Prefix for the created resource names (e.g. \"tf-e2e-customerio\")."
  type        = string
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
