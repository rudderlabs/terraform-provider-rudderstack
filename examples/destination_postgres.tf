resource "rudderstack_destination_postgres" "example" {
  name = "my-postgres"

  config {
    host                 = "localhost"
    user                 = "postgres"
    password             = "postgres"
    port                 = "5432"
    namespace            = "example"
    database             = "example"
    ssl_mode             = "disable"
    use_rudder_storage   = false

    # verify_ca {
    #   client_key  = "-----BEGIN RSA PRIVATE KEY-----...-----END CERTIFICATE-----"
    #   client_cert = "-----BEGIN RSA PRIVATE KEY-----...-----END CERTIFICATE-----"
    #   server_ca   = "-----BEGIN RSA PRIVATE KEY-----...-----END CERTIFICATE-----"
    # }

    # s3 {
    #   bucket_name   = ""
    #   access_key_id = ""
    #   access_key    = ""
    # }

    # gcs {
    #   bucket_name = ""
    #   credentials = ""
    # }

    # azure_blob {
    #   container_name = ""
    #   account_name   = ""
    #   account_key    = ""
    # }

    # minio {
    #   bucket_name       = ""
    #   endpoint          = ""
    #   access_key_id     = ""
    #   secret_access_key = ""
    #   use_ssl           = ""
    # }
  }

  #   sync {
  #     frequency = "30"
  #     sync_start_at    = "???"
  #     exclude_start_at = "???"
  #     exclude_end_at   = "???"
  #   }
}
