resource "rudderstack_destination_amplitude" "example" {
  name = "my-kinesis-tf"

  config {
    api_key = "usa-east"
#    stream = "test"

#    role_based_authentication {
#      i_am_role_arn = "arn-exp"
#    }

    # access_key_id = ""
    # access_key    = ""
    # i_am_role_arn    = "arm"

    # role_based_auth = true
    # use_message_id   = false

    # onetrust_cookie_categories = ["one", "two", "three"]

  }
}
