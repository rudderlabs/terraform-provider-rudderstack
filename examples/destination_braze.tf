resource "rudderstack_destination_braze" "example" {
  name = "my-braze"

  config {
    data_center = "US-01"
    connection_mode {
      web = "cloud"
    }
    rest_api_key    = "braze rest api key"
    # app_key    = "braze app key"
    enable_subscription_group_in_group_call    = true
    enable_nested_array_operations    = true
    send_purchase_event_with_extra_properties    = true
    track_anonymous_user    = true
    support_dedup    = true
    enable_braze_logging {
      web = true
    }
    enable_push_notification {
      web = true
    }
    allow_user_supplied_javascript {
      web = true
    }
    event_filtering {
       whitelist = ["one", "two", "three"]
       # blacklist = ["one", "two", "three"]
    }
    # onetrust_cookie_categories {
    #   web = ["one", "two", "three"]
    #   android = ["one", "two", "three"]
    #   ios = ["one", "two", "three"]
    #   unity = ["one", "two", "three"]
    #   reactnative = ["one", "two", "three"]
    #   flutter = ["one", "two", "three"]
    #   cordova = ["one", "two", "three"]
    #   amp = ["one", "two", "three"]
    #   cloud = ["one", "two", "three"]
    #   warehouse = ["one", "two", "three"]
    #   shopify = ["one", "two", "three"]
    # }
  }
}
