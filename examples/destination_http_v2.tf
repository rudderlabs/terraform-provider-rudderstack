resource "rudderstack_destination_http_v2" "example" {
  name = "my-http-jsonata-destination"

  config {
    api_url = "https://example.com/webhooks/events"

    auth          = "apiKeyAuth"
    api_key_name  = "x-api-key"
    api_key_value = "your-api-key"

    method = "POST"
    format = "JSON"

    request_body_template  = <<-JSONATA
      [$$.events.{"customer": {"id": userId}, "items": properties.items}]
    JSONATA
    supported_message_types = ["track", "identify"]

    is_batching_enabled = true
    max_batch_size      = "5000"
  }
}
