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
    # onetrust_cookie_categories {
    #   web = ["one", "two", "three"]
    #   android = ["one", "two", "three"]
    #   ios = ["one", "two", "three"]
    #   unity = ["one", "two", "three"]
    #   reactnative = ["one", "two", "three"]
    #   flutter = ["one", "two", "three"]
    #   cordova = ["one", "two", "three"]
    #   amp = ["one", "two", "three"]
    #   cloud = ["one", "two", "three"]
    #   warehouse = ["one", "two", "three"]
    #   shopify = ["one", "two", "three"]
    # }
  }
}
