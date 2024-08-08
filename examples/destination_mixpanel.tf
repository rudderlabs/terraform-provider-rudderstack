resource "rudderstack_destination_mixpanel" example{
  name = "my-mixpanel"

  config {
    token = "..."
    data_residency = "us"
    persistence = "none"
    # api_secret = "..."
    # people = true
    # set_all_traits_by_default = true
    # consolidated_page_calls = true
    # track_categorized_pages = true
    # track_named_pages = true
    # source_name = "my-mixpanel"
    # cross_subdomain_cookie = true
    # secure_cookie = true
    # super_properties = ["one","two","three"]
    # people_properties = ["one","two","three"]
    # event_increments = ["one","two","three"]
    # prop_increments = ["one","two","three"]
    # group_key_settings = ["one","two","three"]
    
    # use_native_sdk {
    #   web = true
    # }

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one","two","three"]
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
    # use_new_mapping = true
  }
}
