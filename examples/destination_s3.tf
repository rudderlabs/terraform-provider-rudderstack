resource "rudderstack_destination_s3" "example" {
  name = "my-s3"

  config {
    bucket_name = "bucket"

    # prefix        = "prefix"
    # access_key_id = "..."
    # access_key    = "..."

    # enable_sse    = true
    # onetrust_cookie_categories = ["one", "two", "three"]
  }
}
