resource "rudderstack_destination_deltalake" "example" {
  name = "my-deltalake"

  config {
    host = "abc-12345de-6789.cloud.databricks.com"
    port = "443"
    path = "/sql/1.0/endpoints/..."
    # Personal access token auth (default):
    token = "..."
    # Machine-to-machine OAuth auth instead of a token:
    # use_oauth = true
    # oauth_client_id = "..."
    # oauth_client_secret = "..."
    # catalog = "..."
    # namespace = "..."
    sync {
      frequency = "30"
      # start_at                  = "10:00"
      # exclude_window_start_time = "11:00"
      # exclude_window_end_time   = "12:00"
    }
    # skip_tracks_table = false
    # skip_users_table = true
    # prefer_append = true
    # Create the Delta tables outside the Databricks-managed location:
    # enable_external_location = true
    # external_location = "/mnt/example/delta"
    # Changing either of these replaces the destination:
    # underscore_divide_numbers = false
    # allow_users_context_traits = false

    # Staging object storage. Must match the staging block below.
    bucket_provider = "S3"
    # use_rudder_storage = true  # use RudderStack-managed buckets and omit the staging block
    # cleanup_object_storage_files = false
    # prefix = "..."
    # S3 with access keys:
    s3 {
      bucket_name    = "..."
      use_sts_tokens = true
      access_key_id  = "..."
      access_key     = "..."
      # enable_sse = true
    }
    # S3 with IAM role-based auth (conflicts with access_key_id/access_key):
    # s3 {
    #   bucket_name = "..."
    #   use_sts_tokens = true
    #   role_based_authentication {
    #     i_am_role_arn = "arn:aws:iam::123456789012:role/MyRole"
    #   }
    # }
    # gcs {
    #   bucket_name = "..."
    #   credentials = "..."
    # }
    # azure_blob {
    #   container_name = "..."
    #   account_name = "..."
    #   account_key = "..."
    #   # SAS token instead of the account key:
    #   # use_sas_tokens = true
    #   # sas_token = "..."
    #   # enable_hierarchical_namespace = false
    # }
    # connection_mode {
    #   web = "cloud"
    #   android = "cloud"
    #   android_kotlin = "cloud"
    #   ios = "cloud"
    #   ios_swift = "cloud"
    #   unity = "cloud"
    #   amp = "cloud"
    #   cloud = "cloud"
    #   cloud_source = "cloud"
    #   reactnative = "cloud"
    #   flutter = "cloud"
    #   cordova = "cloud"
    #   shopify = "cloud"
    # }
    # consent_management {
    # 	web = [
    # 		{
    # 			provider = "oneTrust"
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
    # 		resolution_strategy = "or"
    # 		consents = ["one_ios", "two_ios", "three_ios"]
    # 	}]
    # }
  }
}
