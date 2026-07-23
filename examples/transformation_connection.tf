resource "rudderstack_transformation_connection" "example" {
  transformation_id = "2C8Vk2wj8qkofy00YzJbvJOGeqa"
  destination_id    = rudderstack_destination_redshift.example.id
}
