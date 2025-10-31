resource "rudderstack_destination_firebase" "example" {
  name = "my-firebase"

  config {
    connection_mode {
      android = "device"
      ios = "device"
    }
    event_filtering {
       whitelist = ["one", "two", "three"]
       # blacklist = ["one", "two", "three"]
    }
  }
}
