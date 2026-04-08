resource "rudderstack_destination_hs" "example" {
  name = "my-hubspot"

  config {
    authorization_type = "newPrivateAppApi"
    api_version        = "newApi"
    access_token       = "demo_access_token"
    hub_id             = "74X991"
    lookup_field       = "email"
    do_association     = false

    hubspot_events = [
      {
        rs_event_name      = "Product Searched"
        hubspot_event_name = "pe12345678_search"
        event_properties = [
          {
            from = "query"
            to   = "hs_search_query"
          }
        ]
      }
    ]

    connection_mode {
      web = "device"
    }

    use_native_sdk {
      web = true
    }
  }
}
