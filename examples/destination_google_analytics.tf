resource "rudderstack_destination_google_analytics" "example" {
  name = "my-google-analytics"

  config {
    tracking_id = "UA-0000-0000"

    # double_click              = true
    # enhanced_link_attribution = true
    # include_search            = true
    # disable_md5               = true
    # anonymize_ip              = true
    # enhanced_ecommerce        = true
    # non_interaction           = true

    # server_side_identify {
    #   event_category = "..."
    #   event_action   = "..."
    # }

    # track_categorized_pages {
    #   web = true
    # }

    # track_named_pages {
    #   web = true
    # }

    # sample_rate {
    #   web = "1000"
    # }

    # site_speed_sample_rate {
    #   web = "1000"
    # }

    # set_all_mapped_props {
    #   web = true
    # }

    # domain {
    #   web = "..."
    # }

    # optimize {
    #   web = "..."
    # }

    # use_google_amp_client_id {
    #   web = true
    # }

    # named_tracker {
    #   web = true
    # }

    # use_native_sdk {
    #   web = true
    # }

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one", "two", "three"]
    # }

    # reset_custom_dimensions_on_page {
    #   web = ["one", "two", "three"]
    # }

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

    # content_groupings = [{
    #   from = "from"
    #   to   = "to"
    # }]

    # dimensions = [{
    #   from = "from"
    #   to   = "to"
    # }]
  }
}
