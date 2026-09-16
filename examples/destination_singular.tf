resource "rudderstack_destination_singular" "example" {
  name = "my-singular"

  config {
    api_key = "your_api_key_here"

    connection_mode {
      android = "device"
      ios     = "device"
      web     = "cloud"
    }
  }
}
