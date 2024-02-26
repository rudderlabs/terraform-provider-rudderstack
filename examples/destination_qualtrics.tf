resource "rudderstack_destination_qualtrics" example{
  name = "my-qualtrics"

  config {
    project_id = "p-id"
    brand_id = "b-id"
    # enable_generic_page_title = true
    
    # use_native_sdk {
    #   web = true
    # }

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one","two","three"]
    # }

    # onetrust_cookie_categories = ["one", "two", "three"]
  }
}