resource "rudderstack_destination_sentry" "example" {
  name = "my-sentry"

  config {
    dsn = "https://example.slack.com"

    # server_name             = "..."
    # release                 = "..."
    # environment             = "..."
    # custom_version_property = "..."
    # logger                  = "..."
    # debug_mode              = false

    # ignore_errors = ["one", "two", "three"]
    # include_paths = ["one", "two", "three"]
    # allow_urls    = ["one", "two", "three"]
    # deny_urls     = ["one", "two", "three"]

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
