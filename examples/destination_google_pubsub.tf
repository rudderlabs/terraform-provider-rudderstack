resource "rudderstack_destination_google_pubsub" "example" {
  name = "my-google-pubsub-test"

  config {
    project_id = "project-id"
    credentials = "..."

#    event_to_topic_map = [
#      {
#        from = "event-1"
#        to   = "topic-1"
#      }
#    ]

#    event_to_attribute_map = [
#      {
#        from = "event-1"
#        to   = "attribute-1"
#      }
#    ]

#     consent_management {
#     	web = [
#     		{
#     			provider = "oneTrust"
#     			consents = ["one_web", "two_web", "three_web"]
#     			resolution_strategy = ""
#     		},
#     		{
#     			provider = "ketch"
#     			consents = ["one_web", "two_web", "three_web"]
#     			resolution_strategy = ""
#     		},
#     		{
#     			provider = "custom"
#     			resolution_strategy = "and"
#     			consents = ["one_web", "two_web", "three_web"]
#     		}
#     	]
#     	android = [{
#     		provider = "ketch"
#     		consents = ["one_android", "two_android", "three_android"]
#     		resolution_strategy = ""
#     	}]
#     	ios = [{
#     		provider = "custom"
#     		resolution_strategy = "and"
#     		consents = ["one_ios", "two_ios", "three_ios"]
#     	}]
#     	unity = [{
#     		provider = "custom"
#     		resolution_strategy = "or"
#     		consents = ["one_unity", "two_unity", "three_unity"]
#     	}]
#     	reactnative = [{
#     		provider = "custom"
#     		resolution_strategy = "and"
#     		consents = ["one_reactnative", "two_reactnative", "three_reactnative"]
#     	}]
#     	flutter = [{
#     		provider = "custom"
#     		resolution_strategy = "and"
#     		consents = ["one_flutter", "two_flutter", "three_flutter"]
#     	}]
#     	cordova = [{
#     		provider = "custom"
#     		resolution_strategy = "and"
#     		consents = ["one_cordova", "two_cordova", "three_cordova"]
#     	}]
#     	amp = [{
#     		provider = "custom"
#     		resolution_strategy = "and"
#     		consents = ["one_amp", "two_amp", "three_amp"]
#     	}]
#     	cloud = [{
#     		provider = "custom"
#     		resolution_strategy = "and"
#     		consents = ["one_cloud", "two_cloud", "three_cloud"]
#     	}]
#     	warehouse = [{
#     		provider = "custom"
#     		resolution_strategy = "and"
#     		consents = ["one_warehouse", "two_warehouse", "three_warehouse"]
#     	}]
#     	shopify = [{
#     		provider = "custom"
#     		resolution_strategy = "and"
#     		consents = ["one_shopify", "two_shopify", "three_shopify"]
#     	}]
#     }

  }
}
