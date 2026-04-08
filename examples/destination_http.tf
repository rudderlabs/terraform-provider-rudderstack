resource "rudderstack_destination_http" "example" {
  name = "my-http-webhook"

  config {
    api_url = "https://example.com/webhooks/events"

    auth          = "apiKeyAuth"
    api_key_name  = "x-api-key"
    api_key_value = "your-api-key"

    method = "POST"
    format = "JSON"
  }
}