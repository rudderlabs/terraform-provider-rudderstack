resource "rudderstack_destination_postgres" "example" {
  name = "my-postgres-tf"

  config {
    host        = "host"
    database    = "database"
    user        = "user"
    password    = "..."
    port        = "1234"
    use_rudder_storage = true
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
    #   cloud_source = ["one", "two", "three"]
    #   shopify = ["one", "two", "three"]
    # }
  }
}
