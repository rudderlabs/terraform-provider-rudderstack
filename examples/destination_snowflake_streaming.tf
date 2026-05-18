resource "rudderstack_destination_snowflake_streaming" "example" {
  name = "my-snowflake-streaming"

  config {
    account   = "your_snowflake_account"
    database  = "your_database"
    warehouse = "your_warehouse"
    user      = "your_user"
    private_key = <<EOT
-----BEGIN PRIVATE KEY-----
YOUR_PRIVATE_KEY
-----END PRIVATE KEY-----
EOT
    namespace = "your_namespace"

    # role                    = "your_role"
    # private_key_passphrase  = "your_passphrase"
    # skip_tracks_table       = false
    # json_paths              = "event.properties.key1,event.properties.key2"
    # enable_iceberg          = false
    # external_volume         = "EXTERNAL_VOLUME"
    # underscore_divide_numbers = false
    # allow_users_context_traits = false

    # connection_mode {
    #   web          = "cloud"
    #   android      = "cloud"
    #   android_kotlin = "cloud"
    #   ios          = "cloud"
    #   ios_swift    = "cloud"
    #   reactnative  = "cloud"
    #   cloud        = "cloud"
    #   cloud_source = "cloud"
    # }

    # consent_management {
    #   web = [{
    #     provider = "custom"
    #     resolution_strategy = "and"
    #     consents = ["C0001", "C0002"]
    #   }]
    # }
  }
}
