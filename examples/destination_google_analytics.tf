resource "rudderstack_destination_google_analytics" "example" {
  name = "my-google-analytics"

  config {
    tracking_id = "UA-0000-0000"

    # double_click              = true
    # enhanced_link_attribution = true
    # include_search            = true
    # disable_md5               = true
    # anonymize_ip              = true
    # enhanced_ecommerce        = true
    # non_interaction           = true

    # server_side_identify {
    #   event_category = "..."
    #   event_action   = "..."
    # }

    # track_categorized_pages {
    #   web = true
    # }

    # track_named_pages {
    #   web = true
    # }

    # sample_rate {
    #   web = "1000"
    # }

    # site_speed_sample_rate {
    #   web = "1000"
    # }

    # set_all_mapped_props {
    #   web = true
    # }

    # domain {
    #   web = "..."
    # }

    # optimize {
    #   web = "..."
    # }

    # use_google_amp_client_id {
    #   web = true
    # }

    # named_tracker {
    #   web = true
    # }

    # use_native_sdk {
    #   web = true
    # }

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one", "two", "three"]
    # }

    # reset_custom_dimensions_on_page {
    #   web = ["one", "two", "three"]
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
    # 	android = [{
    # 		provider = "ketch"
    # 		consents = ["one_android", "two_android", "three_android"]
    # 		resolution_strategy = ""
    # 	}]
    # 	ios = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_ios", "two_ios", "three_ios"]
    # 	}]
    # 	unity = [{
    # 		provider = "custom"
    # 		resolution_strategy = "or"
    # 		consents = ["one_unity", "two_unity", "three_unity"]
    # 	}]
    # 	reactnative = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_reactnative", "two_reactnative", "three_reactnative"]
    # 	}]
    # 	flutter = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_flutter", "two_flutter", "three_flutter"]
    # 	}]
    # 	cordova = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_cordova", "two_cordova", "three_cordova"]
    # 	}]
    # 	amp = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_amp", "two_amp", "three_amp"]
    # 	}]
    # 	cloud = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_cloud", "two_cloud", "three_cloud"]
    # 	}]
    # 	warehouse = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_warehouse", "two_warehouse", "three_warehouse"]
    # 	}]
    # 	shopify = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_shopify", "two_shopify", "three_shopify"]
    # 	}]
    # }

    # content_groupings = [{
    #   from = "from"
    #   to   = "to"
    # }]

    # dimensions = [{
    #   from = "from"
    #   to   = "to"
    # }]
  }
}
