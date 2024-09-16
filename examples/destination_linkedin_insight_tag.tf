resource "rudderstack_destination_linkedin_insight_tag" "example" {
  name = "my-linkedin-insight-tag"

  config {
    partner_id = "p-id"
#    event_to_conversion_id_map = [
#      {
#        from = "a1"
#        to   = "b1"
#      },
#      {
#        from = "a2"
#        to   = "b2"
#      }]
#    use_native_sdk {
#      web = true
#    }
#    event_filtering {
#      whitelist = ["one", "two", "three"]
#    }
    # onetrust_cookie_categories {
    #   web = ["one", "two", "three"]
    # }
    }
  
}
