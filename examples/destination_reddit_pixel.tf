resource "rudderstack_destination_reddit_pixel" "example" {
  name = "my-reddit-pixel"

  config {
    advertiser_id = "your-reddit-pixel-id"

    # event_mapping_from_config = [{
    #   from = "Order Completed"
    #   to   = "Purchase"
    # }]

    # event_filtering {
    #   whitelist = ["Product Viewed", "Order Completed"]
    # }

    # use_native_sdk {
    #   web = true
    # }

    # connection_mode {
    #   web = "device"
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
