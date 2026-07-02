output "account_id" {
  description = "ID of the created BigQuery rETL account."
  value       = rudderstack_account_source_bigquery.acct.id
}

output "source_id" {
  description = "ID of the created rETL source table."
  value       = rudderstack_retl_source_table.users.id
}
