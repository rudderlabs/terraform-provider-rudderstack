resource "rudderstack_destination_kafka" "example" {
  name = "my-kafka"

  config {
    host_name       = "example.com"
    port            = "9092"
    topic           = "example-topic"
    ssl_enabled     = true
    # ca_certificate  = "asd"
    # use_sasl        = true
    # sasl_type       = "plain"
    # username        = "foo"
    # password        = "Rudder123"
    # convert_to_avro = true
    # # avro_schema = [{
    # #   schema_id = "foo",
    # #   schema    = "bar"
    # # }]
    # embed_avro_schema_id       = true
    # onetrust_cookie_categories = ["analytics"]
  }

}
