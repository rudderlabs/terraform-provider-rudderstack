resource "rudderstack_destination_braze" "example" {
  name = "my-braze"

  config {
    data_center = "US-01"
    connection_mode {
      web = "cloud"
    }
    # rest_api_key    = "braze rest api key"
    # app_key    = "braze app key"
    # enable_subscription_group_in_group_call    = true
    # enable_nested_array_operations    = true
    # send_purchase_event_with_extra_properties    = true
    # track_anonymous_user    = true
    # support_dedup    = true
    # enable_braze_logging    = true
    # enable_push_notification    = true
    # allow_user_supplied_javascript    = true
    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one", "two", "three"]
    # }
    # onetrust_cookie_categories = ["one", "two", "three"]
  }
}
