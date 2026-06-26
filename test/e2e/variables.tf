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

# Customer.io creds are optional: supply api_key + site_id to also exercise the
# BigQuery→Customer.io chain; leave empty to run the webhook-only smoke.
variable "customerio_api_key" {
  description = "Customer.io App API key. Empty skips the BigQuery→Customer.io chain."
  type        = string
  default     = ""
  sensitive   = true
}

variable "customerio_site_id" {
  description = "Customer.io site ID (required with customerio_api_key to enable the chain)."
  type        = string
  default     = ""
  sensitive   = true
}

variable "customerio_datacenter" {
  description = "Customer.io data center (US or EU)."
  type        = string
  default     = "US"
}

# Customer.io Audience creds are optional and independent of the Customer.io
# (VDM v2) chain above: supply customerio_audience_app_api_key + a real
# customerio_audience_id (alongside customerio_site_id + customerio_api_key) to
# also exercise the BigQuery→Customer.io Audience chain. The Audience
# destination reuses site_id/api_key (same Customer.io account) but additionally
# needs the App API key and a region.
variable "customerio_audience_app_api_key" {
  description = "Customer.io App API key for the Audience destination. Empty skips the BigQuery→Customer.io Audience chain."
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
  description = "Customer.io audience ID (positive integer). 0 skips the BigQuery→Customer.io Audience chain."
  type        = number
  default     = 0
}
