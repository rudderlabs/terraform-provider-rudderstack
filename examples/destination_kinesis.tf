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
    # onetrust_cookie_categories = ["one", "two", "three"]

  }
}
