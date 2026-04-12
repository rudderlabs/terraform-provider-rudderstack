
resource "rudderstack_destination_snowflake" "example" {
  name = "my-snowflake"

  config {
    account = "..."
    database = "..."
    warehouse = "..."
    user = "..."
    # Password-based auth (default):
    password = "..."
    # Key pair auth (set use_key_pair_auth = true to use instead of password):
    # use_key_pair_auth = true
    # private_key = "MIIEvQIBADA..."  # raw base64 key body or full PEM format
    # private_key_passphrase = "..."  # only needed if the private key is encrypted
    sync {
      frequency = "60"
      # start_at                  = "10:00"
      # exclude_window_start_time = "11:00"
      # exclude_window_end_time   = "12:00"
    }
    # json_paths = "..."
    use_rudder_storage = true
    # namespace = "..."
    # prefix = "..."
    # additional_properties = true
    # S3 with access keys:
    # s3 {
    #   bucket_name = "..."
    #   access_key_id = "..."
    #   access_key = "..."
    #   enable_sse = true
    #   storage_integration = "..."
    # }
    # S3 with IAM role-based auth (conflicts with access_key_id/access_key):
    # s3 {
    #   bucket_name = "..."
    #   role_based_authentication {
    #     i_am_role_arn = "arn:aws:iam::123456789012:role/MyRole"
    #   }
    #   storage_integration = "..."
    # }
    # gcp {
    #   bucket_name = "..."
    #   credentials = "..."
    #   storage_integration = "..."
    # }
    # azure {
    #   container_name = "..."
    #   account_name = "..."
    #   account_key = "..."
    #   storage_integration = "..."
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
    # 	cloud_source = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_cloud_source", "two_cloud_source", "three_cloud_source"]
    # 	}]
    # 	shopify = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_shopify", "two_shopify", "three_shopify"]
    # 	}]
    # }
  }
}
