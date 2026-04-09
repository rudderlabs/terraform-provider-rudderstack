resource "rudderstack_destination_linkedin_ads" "example" {
  name = "my-linkedin-ads"

  config {
    rudder_account_id = "your-rudder-account-id"
    hash_data         = true
    ad_account_id     = "123456789"

    # deduplication_key = "messageId"

    conversion_mapping = [
      {
        from = "Order Completed"
        to   = "987654321"
      }
    ]

    # consent_management {
    #   web = [
    #     {
    #       provider            = "oneTrust"
    #       consents            = ["consent_web_1", "consent_web_2"]
    #       resolution_strategy = "or"
    #     }
    #   ]
    #   cloud = [
    #     {
    #       provider            = "oneTrust"
    #       consents            = ["consent_cloud_1"]
    #       resolution_strategy = ""
    #     }
    #   ]
    # }
  }
}
