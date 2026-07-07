resource "rudderstack_destination_confluent_cloud" "example" {
  name = "my-confluent-cloud"

  config {
    bootstrap_server = "test.region.provider.confluent.cloud:9092"
    topic            = "example-topic"
    api_key          = "example-api-key"
    api_secret       = "example-api-secret"
  }
}
