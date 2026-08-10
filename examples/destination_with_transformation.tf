# Attach a published transformation to a destination. The transformation must be
# published before it can be attached, and a destination can have at most one.
resource "rudderstack_destination_http" "example" {
  name = "my-filtered-http-webhook"

  config {
    api_url = "https://example.com/webhooks/events"

    auth          = "apiKeyAuth"
    api_key_name  = "x-api-key"
    api_key_value = "your-api-key"

    method = "POST"
    format = "JSON"
  }

  # ID of a published transformation (e.g. one that filters out events that must
  # not reach this destination). Omit or clear to detach.
  transformation_id = "2abcTRANSFORMATIONid"
}
