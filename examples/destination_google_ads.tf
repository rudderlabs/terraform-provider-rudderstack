resource "rudderstack_destination_google_ads" "example" {
  name = "my-google-ads"

  config {
    conversion_id = "AW-00000000"

    # default_page_conversion = "..."

    # page_load_conversions = [
    #   {
    #     "label" = "..."
    #     "name"  = "..."
    #   }
    # ]

    # click_event_conversions = [
    #   {
    #     "label" = "..."
    #     "name"  = "..."
    #   }
    # ]

    # dynamic_remarketing {
    #   web = true
    # }

    # conversion_linker          = true
    # send_page_view             = true
    # disable_ad_personalization = true

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
