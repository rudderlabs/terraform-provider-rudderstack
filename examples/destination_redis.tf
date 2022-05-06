resource "rudderstack_destination_redis" "example" {
  name = "my-redis"

  config {
    address = "1.2.3.4"

    # password      = ""
    # cluster_mode  = false
    # secure        = false
    # prefix        = "..."
    # database      = "..."
    # ca_certificate = "..."
    # skip_verify   = false
  }
}
