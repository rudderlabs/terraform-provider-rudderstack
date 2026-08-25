resource "rudderstack_destination_google_adwords_offline_conversions" "example" {
  name = "my-google-ads-offline-conversions"

  config {
    rudder_account_id = "your-rudder-account-id"
    customer_id       = "1234567890"

    # sub_account       = true
    # login_customer_id = "0987654321"

    events_to_offline_conversions_type_mapping = [
      {
        from = "Order Completed"
        to   = "click"
      }
    ]

    events_to_conversions_names_mapping = [
      {
        from = "Order Completed"
        to   = "Purchase Conversion"
      }
    ]

    custom_variables = [
      {
        from = "revenue"
        to   = "cart_value"
      }
    ]

    user_identifier_source  = "none"
    conversion_environment  = "none"
    default_user_identifier = "email"
    hash_user_identifier    = true
    validate_only           = false

    connection_mode {
      web = "cloud"
    }

    # consent_management {
    #   web = [
    #     {
    #       provider            = "oneTrust"
    #       consents            = ["consent_category_1", "consent_category_2"]
    #       resolution_strategy = ""
    #     }
    #   ]
    # }
  }
}
