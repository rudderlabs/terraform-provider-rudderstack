resource "rudderstack_destination_iterable" "example" {
  name = "my-iterable"

  config {
      api_key = "73983282843839749873"
      
      # map_to_single_event = true
      # track_all_pages = false
      # track_categorized_pages = true
      # track_named_pages = true
      # use_native_sdk {
      #   web = true
      # }
      # initialisation_identifier { 
      #   web = "email" 
      # }
      # get_in_app_event_mapping {
      #   web = ["one", "two", "three"]
      # }
      # purchase_event_mapping { 
      #   web = ["one", "two", "three"]
      # }
      # send_track_for_inapp { 
      #   web = true 
      # }
      # animation_duration { 
      #   web = "200" 
      # }
      # display_interval { 
      #   web = "2500" 
      # }
      # on_open_screen_reader_message { 
      #   web =  "" 
      # }
      # on_open_node_to_take_focus { 
      #   web =  ""  
      # }
      # package_name { 
      #   web = "my-package-test" 
      # }
      # right_offset { 
      #   web = "15" 
      # }
      # top_offset { 
      #   web = "11" 
      # }
      # bottom_offset { 
      #   web = "24%" 
      # }
      # handle_links { 
      #   web = "open-all-new-tab" 
      # }
      # close_button_color { 
      #   web = "blue" 
      # }
      # close_button_size { 
      #   web = ""
      # }
      # close_button_color_top_offset { 
      #   web = "3%"
      # }
      # close_button_color_side_offset { 
      #   web = "2%" 
      # }
      # icon_path { 
      #   web = "" 
      # }
      # is_required_to_dismiss_message { 
      #   web = true 
      # }
      # close_button_position { 
      #   web = "" 
      # }
      # onetrust_cookie_categories = ["one", "two", "three"]

    }
  
}