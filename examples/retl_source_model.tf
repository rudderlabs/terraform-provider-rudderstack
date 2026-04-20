resource "rudderstack_retl_source_model" "users_revenue" {
  name                   = "users-revenue"
  source_definition_name = "snowflake"
  account_id             = "acc-1234"

  config {
    primary_key = "user_id"
    sql         = "select user_id, sum(amount) as revenue from orders group by user_id"
    description = "Revenue per user, refreshed nightly."
  }
}
