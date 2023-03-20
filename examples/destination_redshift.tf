resource "rudderstack_destination_redshift" "example" {
  name = "my-redshift"

  config {
    host     = "localhost"
    port     = "5432"
    database = "example"
    user     = "redshift"
    password = "redshift"

    namespace          = "example"
    enable_sse         = true
    use_rudder_storage = false


    # s3 {
    #   bucket_name   = ""
    #   access_key_id = ""
    #   access_key    = ""
    # }

    # onetrust_cookie_categories = ["one", "two", "three"]

    sync {
      frequency = "30"

      # start_at                  = "10:00"
      # exclude_window_start_time = "11:00"
      # exclude_window_end_time   = "12:00"
    }
  }
}
