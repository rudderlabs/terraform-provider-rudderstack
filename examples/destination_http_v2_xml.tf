resource "rudderstack_destination_http_v2" "xml" {
  name = "my-http-xml-destination"

  config {
    api_url      = "https://example.com/events.xml"
    format       = "XML"
    xml_root_key = "Events"

    request_body_template = <<-JSONATA
      {"event": $$.events.{"id": messageId}}
    JSONATA
    max_batch_size = "10"
  }
}

resource "rudderstack_destination_http_v2" "form" {
  name = "my-http-form-destination"

  config {
    api_url               = "https://example.com/form"
    format                = "FORM"
    request_body_template = <<-JSONATA
      {"email": $$.events[0].traits.email}
    JSONATA
  }
}
