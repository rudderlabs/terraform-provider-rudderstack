resource "rudderstack_destination_adobe_analytics" "example" {
    name = "my-adobe-analytics"

    config {
        tracking_server_url = "http://sampleurl.com"
        report_suite_id = "id01, id02"
        heartbeat_tracking_server_url = "http://sampleheartbeaturl.com"
    }
}