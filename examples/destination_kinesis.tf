resource "rudderstack_destination_kinesis" "example" {
  name = "my-kinesis-tf"

  config {
    region = "usa-east"
    stream = "test"

    role_based_authentication {
      i_am_role_arn = "arn-exp"
    }

    # key_based_authentication {
          # access_key_id = ""
          # access_key    = ""
    #    }
    # use_message_id   = false
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
