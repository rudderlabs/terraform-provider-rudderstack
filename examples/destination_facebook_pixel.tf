resource "rudderstack_destination_facebook_pixel" "example" {
  name = "my-facebook-pixel"

  config {
    pixel_id = "facebook pixel id"

    # access_token = "facebook access token"

    # standard_page_call     = true
    # value_field_identifier = "properties.price"
    # advanced_mapping       = true
    # test_destination       = true
    # test_event_code        = "..."

    # events_to_events = [{
    #   from = "from"
    #   to   = "to"
    # }]

    # event_custom_properties = ["one", "two", "three"]

		# blacklist_pii_properties = [{ 
		# 	property = "one"
		# 	hash     = false
		# }, { 
		# 	property = "two"
		# 	hash     = true
		# }]

		# whitelist_pii_properties = [{ 
		# 	property = "one"
		# }, { 
		# 	property = "two"
		# }]
    
    # category_to_content = [{
    #   from = "from"
    #   to   = "to"
    # }]

    # legacy_conversion_pixel_id {
    #   from = "from"
    #   to   = "to"
    # }

    # use_native_sdk {
    #   web = true
    # }

    # event_filtering {
    #   whitelist = ["one", "two", "three"]
    #   blacklist = ["one", "two", "three"]
    # }

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
