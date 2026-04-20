resource "rudderstack_retl_source_table" "users" {
  name                   = "users"
  source_definition_name = "snowflake"
  account_id             = "acc-1234"

  config {
    primary_key = "id"
    schema      = "public"
    table       = "users"
  }
}
