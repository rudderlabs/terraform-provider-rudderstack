resource "rudderstack_destination_reddit" "example" {
  name = "my-reddit"

  config {
    rudder_account_id = "your-rudder-account-id"
    account_id        = "your-reddit-pixel-id"

    # version   = "v3"
    # hash_data = true

    # events_mapping = [{
    #   from = "Order Completed"
    #   to   = "Purchase"
    # }]

    # connection_mode {
    #   web           = "cloud"
    #   android_kotlin = "cloud"
    #   ios_swift      = "cloud"
    #   reactnative    = "cloud"
    # }

    # consent_management {
    #   web = [{
    #     provider            = "oneTrust"
    #     consents            = ["category1", "category2"]
    #     resolution_strategy = ""
    #   }]
    # }
  }
}
