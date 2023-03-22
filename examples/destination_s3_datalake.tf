resource "rudderstack_destination_s3_datalake" "example" {
  name = "my-s3-datalake"

  config {
    bucket_name = "bucket"

    # namespace     = "..."
    # prefix        = "prefix"
    # access_key_id = "..."
    # access_key    = "..."

    # enable_sse    = true

    # onetrust_cookie_categories = ["one", "two", "three"]

    use_glue = true
    region   = "us-east-2"

    sync {
      frequency = "30"
      # start_at  = "10:00"
    }
  }
}
