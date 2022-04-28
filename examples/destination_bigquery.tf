resource "rudderstack_destination_bigquery" "example" {
  name = "my-bigquery"

  config {
    project     = "project"
    bucket_name = "bucket"
    credentials = "..."

    # location  = "us-east1"
    # prefix    = ""
    # namespace = ""

    sync {
      frequency = "30"
      # start_at                  = "10:00"
      # exclude_window_start_time = "11:00"
      # exclude_window_end_time   = "12:00"
    }
  }
}
