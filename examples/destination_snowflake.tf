
resource "rudderstack_destination_snowflake" example{
  name = "my-snowflake"

  config {
    account = "..."
    database = "..."
    warehouse = "..."
    user = "..."
    password = "..."
    sync {
      frequency = "60"
      # start_at                  = "10:00"
      # exclude_window_start_time = "11:00"
      # exclude_window_end_time   = "12:00"
    }
    # json_paths = "..."
    use_rudder_storage = true
    # role = "..."
    # namespace = "..."
    # prefix = "..."
    # additional_properties = true
    # s3 {
    #   bucket_name = "..."
    #   access_key_id = "..."
    #   access_key = "..."
    #   enable_sse = true
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
    # onetrust_cookie_categories = ["one", "two", "three"]
  }
}
