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

    # consent_management {
    # 	web = [
    # 		{
    # 			provider = "oneTrust"
    # 			consents = ["one_web", "two_web", "three_web"]
    # 			resolution_strategy = ""
    # 		},
    # 		{
    # 			provider = "ketch"
    # 			consents = ["one_web", "two_web", "three_web"]
    # 			resolution_strategy = ""
    # 		},
    # 		{
    # 			provider = "custom"
    # 			resolution_strategy = "and"
    # 			consents = ["one_web", "two_web", "three_web"]
    # 		}
    # 	]
    # }
  }
}
