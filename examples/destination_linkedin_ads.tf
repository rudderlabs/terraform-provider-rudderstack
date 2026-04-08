resource "rudderstack_destination_linkedin_ads" "example" {
  name = "my-linkedin-ads"

  config {
    rudder_account_id = "your-rudder-account-id"
    hash_data         = true
    ad_account_id     = "123456789"

    conversion_mapping = [
      {
        from = "Order Completed"
        to   = "987654321"
      }
    ]
  }
}
