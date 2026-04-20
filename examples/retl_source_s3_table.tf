resource "rudderstack_retl_source_s3_table" "events" {
  name       = "s3-events"
  account_id = "acc-1234"

  config {
    bucket_name   = "my-events-bucket"
    object_prefix = "events/2024/"
  }
}
