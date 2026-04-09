resource "rudderstack_destination_customerio" "example" {
  name = "my-customerio"

  config {
    site_id = "customer io site id"
    api_key = "customer io api key"

    # device_token_event_name = ""

    # datacenter = "US"

    # use_native_sdk {
    #   web     = true
    #   android = true
    #   ios     = true
    # }

    # send_page_name_in_sdk {
    #   web = true
    # }

    # data_use_in_app {
    #   web = false
    # }

    # auto_track_device_attributes {
    #   android = true
    #   ios     = true
    # }

    # background_queue_min_number_of_tasks {
    #   android = "10"
    # }

    # background_queue_seconds_delay {
    #   android = "30"
    # }

    # event_filtering {
    #   whitelist = ["event-one", "event-two"]
    #   # blacklist = ["event-one", "event-two"]  # use either whitelist or blacklist, not both
    # }

    # consent_management {
    #   web = [
    #     {
    #       provider            = "oneTrust"
    #       consents            = ["category-1", "category-2"]
    #       resolution_strategy = ""
    #     },
    #     {
    #       provider            = "custom"
    #       resolution_strategy = "and"
    #       consents            = ["category-1", "category-2"]
    #     }
    #   ]
    #   android = [{
    #     provider            = "oneTrust"
    #     consents            = ["category-1", "category-2"]
    #     resolution_strategy = ""
    #   }]
    #   ios = [{
    #     provider            = "oneTrust"
    #     consents            = ["category-1", "category-2"]
    #     resolution_strategy = ""
    #   }]
    # }
  }
}
