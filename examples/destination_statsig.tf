resource "rudderstack_destination_statsig" "example" {
  name = "my-statsig-tf"

  config {
    secret_key = "key"
    connection_mode {
      web = "cloud"
    }
    # onetrust_cookie_categories = ["one", "two", "three"]
  }
}
