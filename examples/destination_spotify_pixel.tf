resource "rudderstack_destination_spotify_pixel" "example" {
  name = "my-spotify-pixel"

  config {
    pixel_id = "your-spotify-pixel-id"

    # enable_alias_call = true

    # events_to_spotify_pixel_events = [{
    #   from = "Order Completed"
    #   to   = "purchase"
    # }]

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    # }

    # connection_mode {
    #   web = "device"
    # }

    # consent_management {
    #   web = [{
    #     provider            = "oneTrust"
    #     consents            = ["category1", "category2"]
    #     resolution_strategy = ""
    #   }]
    # }
  }
}
