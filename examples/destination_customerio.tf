resource "rudderstack_destination_customerio" "example" {
  name = "my-customerio"

  config {
    site_id = "customer io site id"
    api_key = "customer io api key"

    # device_token_event_name = ""

    # datacenter_eu = true

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
  }
}
