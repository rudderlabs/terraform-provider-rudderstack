resource "rudderstack_destination_adobe_analytics" "example" {
    name = "my-adobe-analytics"

    config {
        tracking_server_url = "http://sampleurl.com"
        report_suite_id = "id01, id02"
        # heartbeat_tracking_server_url = "http://sampleheartbeaturl.com"
        # events_to_types = [{
		# 	from = "video start"
		# 	to = "heartbeatPlaybackStarted"
		# }]
		# list_delimiter = [{
		# 	from = "listPhone"
		# 	to = ","
		# }]
		# props_delimiter = [{
		# 	from = "customPhone"
		# 	to = ","
		# }]
		# event_merch_properties = [
		# 	"currency"
		# ]
		# product_merch_properties = [
		# 	"currency"
		# ]
		# event_filtering{
		# 	blacklist = ["one", "two", "three"]
		# }
		# rudder_events_to_adobe_events = [{
		# 	from = "product searched"
		# 	to = "ps1,ps2"
		# }]
		# context_data_mapping = [{
		# 	from = "page.name"
		# 	to = "pName"
		# }]
		# mobile_event_mapping = [{
		# 	from = "page.name"
		# 	to = "pName"
		# }]
		# e_var_mapping = [{
		# 	from = "phone"
		# 	to = "1"
		# }]
		# hier_mapping = [{
		# 	from = "phone"
		# 	to = "1"
		# }]
		# list_mapping = [{
		# 	from = "listPhone"
		# 	to = "1"
		# }]
		# custom_props_mapping = [{
		# 	from = "phone"
		# 	to = "1"
		# }]
		# event_merch_event_to_adobe_event = [{
		# 	from = "Order Completed"
		# 	to = "merchEvent1"
		# }]
		# product_merch_event_to_adobe_event = [{
		# 	from = "Product Ordered"
		# 	to = "MerchProduct1"
		# }]
		# product_merch_evars_map = [{
		# 	from = "phone"
		# 	to = "1"
		# }]
    }
}