resource "rudderstack_account_source_postgres" "example" {
  name = "my-postgres-account"
  config {
    host   = "db.example.com"
    dbname = "analytics"
    user   = "rudder"
    port   = "5432"
    # "disable" (default) or "require".
    ssl_mode = "require"
    # Keep the password out of version control (use a variable or a secret store).
    password = var.postgres_password
  }
}
