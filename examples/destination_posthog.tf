resource "rudderstack_destination_posthog" "example" {
  name = "example-posthog"
  config {
    endpoint = "https://app.posthog.com" # Your PostHog instance URL
    api_key = "..."
    use_v2_group = true
    connection_mode {
      web = "cloud"
      cloud = "cloud"
      flutter = "cloud"
    }
    # autocapture {
    #   web = false
    # }
    # disable_session_recording {
    #   web = false
    # }
    # enable_local_storage_persistence {
    #   web = false
    # }
    # person_profiles {
    #   web = "always"
    # }
    # xhr_headers = [
    #   {
    #     key = "custom-header"
    #     value = "header-value"
    #   }
    # ]
    # property_blacklist = [
    #   {
    #     property = "sensitive_data"
    #   }
    # ]
  }
}