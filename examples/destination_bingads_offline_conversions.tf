resource "rudderstack_destination_bingads_offline_conversions" "example" {
  name = "my-bingads-offline-conversions"

  config {
    rudder_account_id   = "your-rudder-account-id"
    customer_account_id = "53212345"
    customer_id         = "343598"
    is_hash_required    = false

    connection_mode {
      warehouse = "cloud"
    }

    # consent_management {
    #   warehouse = [
    #     {
    #       provider            = "custom"
    #       consents            = ["consent_category_1"]
    #       resolution_strategy = "and"
    #     }
    #   ]
    # }
  }
}
