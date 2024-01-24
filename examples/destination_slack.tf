resource "rudderstack_destination_slack" "example" {
  name = "my-slack"

  config {
    webhook_url = "https://example.slack.com"

    identify_template = "..."

    # event_channel_settings = [
    #   {
    #     name     = "..."
    #     channel  = "..."
    #     regex    = true
    #   }
    # ]

    # event_template_settings = [
    #   {
    #     name     = "..."
    #     template = "..."
    #     regex    = true
    #   }
    # ]

    # whitelisted_trait_settings = ["one", "two", "three"]

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
