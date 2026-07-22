resource "rudderstack_destination_firebase" "example" {
  name = "my-firebase"

  config {
    connection_mode {
      android_kotlin = "device"
      android = "device"
      ios_swift = "device"
      ios = "device"
    }
    event_filtering {
       whitelist = ["one", "two", "three"]
       # blacklist = ["one", "two", "three"]
    }
  }
}
