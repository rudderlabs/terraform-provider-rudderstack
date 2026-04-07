resource "rudderstack_destination_facebook_conversions" "example" {
  name = "my-facebook-conversions"

  config {
    dataset_id   = "1234567898765"
    access_token = "my-access-token"

    # action_source      = "website"
    # limited_data_usage = false
    # test_destination   = false
    # test_event_code    = "TEST80569"
    # remove_external_id = false

    # events_to_events = [{
    #   from = "Product Searched"
    #   to   = "Search"
    # }]

    # blacklist_pii_properties = [{
    #   property = "phone"
    #   hash     = false
    # }]

    # whitelist_pii_properties = [{
    #   property = "name"
    # }]

    # consent_management {
    #   web = [
    #     {
    #       provider            = "oneTrust"
    #       consents            = ["one_web", "two_web", "three_web"]
    #       resolution_strategy = ""
    #     }
    #   ]
    # }
  }
}
