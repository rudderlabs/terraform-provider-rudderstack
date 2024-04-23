resource "rudderstack_destination_webhook" "example" {
  name = "my-webhook"

  config {
    webhook_url    = "https://example.com"
    webhook_method = "GET"

    headers = [
      {
        from = "header-1"
        to   = "value-1"
      },
      {
        from = "header-2"
        to   = "value-2"
      }
    ]

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
