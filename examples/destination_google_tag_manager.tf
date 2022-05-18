resource "rudderstack_destination_google_tag_manager" "example" {
  name = "my-google-tag-manager"

  config {
    container_id = "GTM-000000"

    server_url = "https://example.com"

    use_native_sdk {
      web = true
    }

    event_filtering {
      whitelist = ["one", "two", "three"]
      blacklist = ["one", "two", "three"]
    }

    onetrust_cookie_categories {
      web = ["one", "two", "three"]
    }
  }
}
