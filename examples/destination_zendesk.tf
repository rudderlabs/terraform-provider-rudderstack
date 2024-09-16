resource "rudderstack_destination_zendesk" "example" {
  name = "my-zendesk"

  config {
    email     = "test@example.com"
    api_token = "..."
    domain    = "..."

    create_users_as_verified         = false
    send_group_calls_without_user_id = false
    remove_users_from_organization   = false
    search_by_external_id = false
    # onetrust_cookie_categories {
    #   web = ["one", "two", "three"]
    #   android = ["one", "two", "three"]
    #   ios = ["one", "two", "three"]
    #   unity = ["one", "two", "three"]
    #   reactnative = ["one", "two", "three"]
    #   flutter = ["one", "two", "three"]
    #   cordova = ["one", "two", "three"]
    #   amp = ["one", "two", "three"]
    #   cloud = ["one", "two", "three"]
    #   warehouse = ["one", "two", "three"]
    #   shopify = ["one", "two", "three"]
    # }
  }
}
