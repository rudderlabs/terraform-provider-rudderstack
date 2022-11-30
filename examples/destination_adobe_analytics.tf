resource "rudderstack_destination_adobe_analytics" "example" {
  name = "my-adobe-analytics"

  config {
    tracking_server_url = "http://sampleurl.com"
    tracking_server_secure_url = "http://sampleurl2.com"
    # ssl_heartbeat = true
    # heartbeat_tracking_server_url= "http://rftgh.com"
    # use_utf8_charset = false
    # use_secure_server_side = false
    # proxy_normal_url = "http://rftgh65.com"
    # proxy_heartbeat_url = "http://gh.com"
    # marketing_cloud_org_id = "test_234"
    # drop_visitor_id = "false"
    # timestamp_optional_reporting = "true"
    # no_fallback_visitor_id = "false"
    # prefer_visitor_id = "false"
    # track_page_name = "false"
    # context_data_prefix = "ruddertest"
    # use_legacy_link_name = "false"
    # page_name_fallback_tostring = "false"
    # send_false_values = false
    # product_identifier = "sku"
    #  use_native_sdk {
    #   web          = true
    #   ios          = true
    #   android      = false
    #   react_native = true
    # }
    # onetrust_cookie_categories {
    #   web = ["one", "two", "three"]
    # }
  }
}
