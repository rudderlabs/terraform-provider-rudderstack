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

    # consent_management {
    # 	web = [
    # 		{
    # 			provider = "oneTrust"
    # 			consents = ["one_web", "two_web", "three_web"]
    # 			resolution_strategy = ""
    # 		},
    # 		{
    # 			provider = "ketch"
    # 			consents = ["one_web", "two_web", "three_web"]
    # 			resolution_strategy = ""
    # 		},
    # 		{
    # 			provider = "custom"
    # 			resolution_strategy = "and"
    # 			consents = ["one_web", "two_web", "three_web"]
    # 		}
    # 	]
    # }
  }
}
