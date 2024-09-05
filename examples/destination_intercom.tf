resource "rudderstack_destination_intercom" "example" {
  name = "my-intercom-tf"

  config {
    app_id = "app-id"
    api_key = "api-key"
#    use_native_sdk {
#      web = true
#      ios = true
#    }
#    event_filtering {
#      blacklist = [ "one", "two", "three" ]
#    }
#    onetrust_cookie_categories = ["one", "two", "three"]
#    mobile_api_key_android = "and-key"
#    mobile_api_key_ios = "ios-key"
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
