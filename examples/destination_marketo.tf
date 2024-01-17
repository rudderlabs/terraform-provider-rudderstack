resource "rudderstack_destination_marketo" "example" {
  name = "my-marketo-tf"

  config {
    account_id             = "id1"
    client_id              = "cid2"
    client_secret          = "cs"
    track_anonymous_events = true
    create_if_not_exist    = false
    lead_trait_mapping = [
      {
        from = "property1",
        to   = "value1",
      }
    ]
    rudder_events_mapping = [
      {
        event             = "event0",
        marketo_primarykey = "marketoPrimarykey0",
        marketo_activity_id = "marketoActivityId0",
      }
    ]
    custom_activity_property_map = [
      {
        from = "property1"
        to   = "value1"
      }
    ]
    onetrust_cookie_categories = ["one", "two", "three"]
  }
}
