resource "rudderstack_destination_kinesis" "example" {
  name = "my-marketo-tf"

  config {
    account_id             = "id1"
    client_id              = "cid2"
    client_secret          = "cs"
    track_anonymous_events = true
    create_if_not_exist    = false
    # lead_trait_mapping = [
    #   {
    #     from = "property0"
    #     to   = "value0"
    #   }
    # ]
    # rudder_events_mapping = [
    #   {
    #     event             = "event0"
    #     marketoPrimarykey = "marketoPrimarykey0"
    #     marketoActivityId = "marketoActivityId0"
    #   }
    # ]
    # custom_activity_property_map = [
    #   {
    #     from = "property1"
    #     to   = "value1"
    #   }
    # ]
    # oneTrustCookieCategories = ["one", "two", "three"]
  }
}
