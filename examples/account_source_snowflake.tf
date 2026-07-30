# Key-pair authentication (authentication_type = "keyPair"): supply private_key
# (and optionally private_key_passphrase). Do NOT set password.
resource "rudderstack_account_source_snowflake" "keypair" {
  name = "my-snowflake-keypair-account"
  config {
    account             = "xy12345.eu-west-1"
    dbname              = "ANALYTICS"
    warehouse           = "COMPUTE_WH"
    user                = "RUDDER"
    role                = "ANALYST"
    authentication_type = "keyPair"
    # Load the PEM key from a file kept out of version control.
    private_key            = file("${path.module}/rsa_key.p8")
    private_key_passphrase = var.snowflake_key_passphrase
  }
}

# Password authentication (authentication_type = "password"): supply password.
# Do NOT set private_key / private_key_passphrase.
resource "rudderstack_account_source_snowflake" "password" {
  name = "my-snowflake-password-account"
  config {
    account             = "xy12345.eu-west-1"
    dbname              = "ANALYTICS"
    warehouse           = "COMPUTE_WH"
    user                = "RUDDER"
    authentication_type = "password"
    password            = var.snowflake_password
  }
}
