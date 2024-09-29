resource "rudderstack_destination_amplitude" "example" {
  name = "my-amplitude"

  config {
    api_key    = "amplitude api key"
    api_secret = "amplitude api secret"

    # group_type_trait  = ""
    # group_value_trait = ""

    # track_all_pages           = true
    # track_categorized_pages   = true
    # track_named_pages         = true
    # track_products_once       = true
    # track_revenue_per_product = true

    # track_gclid {
    #   web = true
    # }

    # track_referrer {
    #   web = true
    # }

    # track_utm_properties {
    #   web = true
    # }

    # track_session_events {
    #   android      = true
    #   ios          = true
    #   react_native = true
    # }

    # version_name = ""

    # traits_to_increment = ["one", "two", "three"]
    # traits_to_set_once  = ["one", "two", "three"]
    # traits_to_append    = ["one", "two", "three"]
    # traits_to_prepend   = ["one", "two", "three"]

    # prefer_anonymous_id_for_device_id {
    #   web = true
    # }

    # device_id_from_url_param {
    #   web = true
    # }

    # force_https {
    #   web = true
    # }

    # save_params_referrer_once_per_session {
    #   web = true
    # }

    # unset_params_referrer_on_new_session {
    #   web = true
    # }

    # batch_events {
    #   web = true
    # }

    # map_device_brand = true

    # event_upload_period_millis {
    #   web          = "1000"
    #   ios          = "1000"
    #   android      = "1000"
    #   react_native = "1000"
    # }

    # event_upload_threshold {
    #   web          = "1000"
    #   ios          = "1000"
    #   android      = "1000"
    #   react_native = "1000"
    # }

    # enable_location_listening {
    #   android      = true
    #   react_native = true
    # }

    # use_advertising_id_for_device_id {
    #   android      = true
    #   react_native = true
    # }

    # use_idfa_as_device_id {
    #   ios          = true
    #   react_native = true
    # }

    # use_native_sdk {
    #   web          = true
    #   ios          = true
    #   android      = true
    #   react_native = true
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

    # residency_server = "EU"
  }
}
