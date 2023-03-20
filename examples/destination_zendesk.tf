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
    # onetrust_cookie_categories = ["one", "two", "three"]
  }
}
