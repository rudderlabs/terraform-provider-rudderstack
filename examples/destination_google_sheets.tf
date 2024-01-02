resource "rudderstack_destination_google_sheets" "example" {
  name = "my-google-sheets"

  config {
    sheet_name = "sheet"
    credentials = "..."
    sheet_id = "123"

#    event_key_map = [
#      {
#        from = "header-1"
#        to   = "value-1"
#      },
#      {
#        from = "header-2"
#        to   = "value-2"
#      }
#    ]

    # onetrust_cookie_categories = ["one", "two", "three"]
  }
}
