resource "rudderstack_destination_bigquery" "example" {
  name = "my-bigquery"

  config {
    project     = "project"
    bucket_name = "bucket"
    credentials = "..."

    # location  = "us-east1"
    # prefix    = ""
    # namespace = ""
    # onetrust_cookie_categories {
    #   web = ["one", "two", "three"]
    #   android = ["one", "two", "three"]
    #   ios = ["one", "two", "three"]
    #   unity = ["one", "two", "three"]
    #   reactnative = ["one", "two", "three"]
    #   flutter = ["one", "two", "three"]
    #   cordova = ["one", "two", "three"]
    #   amp = ["one", "two", "three"]
    #   cloud = ["one", "two", "three"]
    #   cloud_source = ["one", "two", "three"]
    #   shopify = ["one", "two", "three"]
    # }

    sync {
      frequency = "30"
      # start_at                  = "10:00"
      # exclude_window_start_time = "11:00"
      # exclude_window_end_time   = "12:00"
    }
  }
}
