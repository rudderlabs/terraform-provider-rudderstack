resource "rudderstack_destination_webhook" "example" {
  name = "my-webhook"

  config {
    webhook_url    = "https://example.com"
    webhook_method = "GET"

    headers = [
      {
        from = "header-1"
        to   = "value-1"
      },
      {
        from = "header-2"
        to   = "value-2"
      }
    ]
  }
}
