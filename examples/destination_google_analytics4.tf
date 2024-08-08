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
