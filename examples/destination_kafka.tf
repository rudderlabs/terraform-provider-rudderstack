resource "rudderstack_destination_kafka" "example" {
  name = "my-kafka"

  config {
    host_name = "example.com"
    port      = "9092"
    topic     = "example-topic"
    # sasl_enabled = true
    # onetrust_cookie_categories = ["analytics"]
  }

}
