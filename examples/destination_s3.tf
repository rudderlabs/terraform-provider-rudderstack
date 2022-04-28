resource "rudderstack_destination_s3" "example" {
  name = "my-s3"

  config {
    bucket_name = "bucket"

    # prefix        = "prefix"
    # access_key_id = "..."
    # access_key    = "..."

    # enable_sse    = true
  }
}
