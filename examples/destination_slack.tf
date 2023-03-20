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

    # onetrust_cookie_categories = ["one", "two", "three"]
  }
}
