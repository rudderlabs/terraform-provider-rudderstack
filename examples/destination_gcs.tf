resource "rudderstack_destination_gcs" "example" {
  name = "my-gcs"

  config {
    bucket_name = "bucket"

    # prefix        = "prefix"
    # credentials   = "..."
  }
}
