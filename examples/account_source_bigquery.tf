resource "rudderstack_account_source_bigquery" "example" {
  name = "my-bigquery-account"
  config {
    project  = "my-gcp-project"
    location = "US"
    # Full GCP service-account key JSON. Load it from a file (keep the file out
    # of version control) rather than inlining the multi-line key.
    credentials = file("${path.module}/service-account-key.json")
  }
}
