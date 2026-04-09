resource "rudderstack_destination_customerio_audience" "example" {
  name = "my-customerio-audience"

  config {
    site_id     = "your-site-id"
    api_key     = "your-api-key"
    app_api_key = "your-app-api-key"
    region      = "US"

    connection_mode {
      warehouse = "cloud"
    }
  }
}
