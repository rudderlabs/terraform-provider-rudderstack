resource "rudderstack_connection" "example" {
  source_id      = rudderstack_source_javascript.example.id
  destination_id = rudderstack_destination_postgres.example.id
}
