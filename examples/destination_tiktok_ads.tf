resource "rudderstack_destination_tiktok_ads" "example" {
  name = "my-tiktok-ads"

  config {
    pixel_code = "A1T8T4XXXXVIQA8ORZMX9"
    access_token = "your-access-token"
  }
}
