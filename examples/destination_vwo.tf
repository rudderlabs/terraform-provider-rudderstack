resource "rudderstack_destination_vwo" "example" {
  name = "my-vwo"

  config {
    account_id = "..."

    # is_spa                   = false
    # send_experiment_track    = false
    # send_experiment_identify = false

    # library_tolerance  = "2000"
    # settings_tolerance = "2000"

    # use_existing_jquery = false

    # use_native_sdk {
    #   web = true
    # }

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one", "two", "three"]
    # }
  }
}
