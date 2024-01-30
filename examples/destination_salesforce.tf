resource "rudderstack_destination_salesforce" "example" {
  name = "my-salesforce"

  config {
    user_name = "user"
    password = "pwd"
    initial_access_token = "token"

    # use_contact_id    = true
    # map_properties    = true
    # onetrust_cookie_categories = ["one", "two", "three"]
  }
}
