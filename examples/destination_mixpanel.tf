resource "rudderstack_destination_mixpanel" example{
  name = "my-mixpanel"

  config {
    token = "avasaffav1241"
    data_residency = "us"
    persistence = "none"
    # api_secret = ""
    people = true
    set_all_traits_by_default = true
    consolidated_page_calls = true
    track_categorized_pages = true
    track_named_pages = true
    source_name = "my-mixpanel"
    cross_subdomain_cookie = true
    secure_cookie = true
    super_properties = [
      {
        "property" = ""
      }
    ]
    people_properties = [
      {
        "property" = ""
      }
    ]
    event_increments = [
      {
        "property" = ""
      }
    ]
    prop_increments = [
      {
        "property" = ""
      }
    ]
    group_key_settings = [
      {
        "group_key" = ""
      }
    ]
    # use_native_sdk = [
    #   {
    #     "web" = true
    #   }
    # ]

    event_filtering {
      whitelist = ["one", "two", "three"]
    }

    onetrust_cookie_categories {
      web = ["one", "two", "three"]
    }
    
  }
}