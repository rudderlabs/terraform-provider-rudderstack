resource "rudderstack_destination_active_campaign" "example" {
  name = "my-active_campaign"

  config {
    api_url = "https://example.api-us1.com"
    api_key = "..."

    # actid     = "..."
    # event_key = "..."
  }
}
