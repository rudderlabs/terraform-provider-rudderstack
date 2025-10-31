resource "rudderstack_destination_adjust" "example" {
  name = "my-adjust"

  config {
    app_token = "your_app_token_here"
    # delay = "5"
    environment = true
    connection_mode {
      android = "device"
      ios     = "device"
    }
    # custom_mappings = [
    #   {
    #     from = "Product Purchased"
    #     to   = "abc123"
    #   },
    #   {
    #     from = "Signup"
    #     to   = "def456"
    #   }
    # ]
    # partner_param_keys = [
    #   {
    #     from = "userId"
    #     to   = "user_id"
    #   }
    # ]
    # enable_install_attribution_tracking {
    #   android = true
    #   ios     = true
    # }
    # event_filtering {
    #    whitelist = ["one", "two", "three"]
    #    # blacklist = ["one", "two", "three"]
    # }
    # consent_management {
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
