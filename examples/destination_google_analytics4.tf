resource "rudderstack_destination_google_analytics4" "example" {
  name = "my-google-analytics4"

  config {
    measurement_id  = "..."

    # firebase_app_id = "..."

    # api_secret = "..."

    # types_of_client = "gtag"

    # block_page_view_event   = false
    # extend_page_view_params = false
    # send_user_id            = false

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
