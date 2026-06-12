resource "rudderstack_account_source_bigquery" "example" {
  name = "my-bigquery-account"
  config {
    project     = "my-gcp-project"
    location    = "US"
    credentials = "{\"type\":\"service_account\"}"
  }
}
