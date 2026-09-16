resource "rudderstack_destination_impact" "example" {
  name = "my-impact"

  config {
    account_sid = "your_account_sid_here"
    api_key     = "your_api_key_here"
    campaign_id = "12345"

    connection_mode {
      web     = "cloud"
      android = "cloud"
      ios     = "cloud"
    }
  }
}
