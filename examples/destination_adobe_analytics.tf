resource "rudderstack_destination_adobe_analytics" "example" {
    name = "my-adobe-analytics"

    config {
        tracking_server_url = ""
        report_suite_id = ""
        ss_heartbeat = ""
        heartbeat_tracking_server_url = ""
        events_to_types = ""
        marketing_cloud_org_id = ""
    }
}