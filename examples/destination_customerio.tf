resource "rudderstack_destination_customerio" "example" {
  name = "my-customerio"

  config {
    site_id = "customer io site id"
    api_key = "customer io api key"

    # device_token_event_name = ""

    # datacenter_eu = true

    # use_native_sdk {
    #   web = true
    # }

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one", "two", "three"]
    # }

    # onetrust_cookie_categories = ["one", "two", "three"]
  }
}
