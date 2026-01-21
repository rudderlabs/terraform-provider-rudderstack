resource "rudderstack_destination_bqstream" "example" {
  name = "my-bqstream"

  config {
    projectId = "your_project_id_here"
    datasetId = "your_dataset_id_here"
    tableId = "your_table_id_here"
    # insertId = "your_insert_id_here"
    credentials = "your_credentials_here"
    connection_mode {
      web = "cloud"
    }
    # consent_management {
    # 	android = [{
    # 		provider = "ketch"
    # 		consents = ["one_android", "two_android", "three_android"]
    # 		resolution_strategy = ""
    # 	}]
    # 	androidKotlin = [{
    # 		provider = "ketch"
    # 		consents = ["one_androidKotlin", "two_androidKotlin", "three_androidKotlin"]
    # 		resolution_strategy = ""
    # 	}]
    # 	ios = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_ios", "two_ios", "three_ios"]
    # 	}]
    # 	iosSwift = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_iosSwift", "two_iosSwift", "three_iosSwift"]
    # 	}]
    # 	unity = [{
    # 		provider = "custom"
    # 		resolution_strategy = "or"
    # 		consents = ["one_unity", "two_unity", "three_unity"]
    # 	}]
    # 	amp = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_amp", "two_amp", "three_amp"]
    # 	}]
    # 	reactnative = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_reactnative", "two_reactnative", "three_reactnative"]
    # 	}]
    # 	flutter = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_flutter", "two_flutter", "three_flutter"]
    # 	}]
    # 	cordova = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_cordova", "two_cordova", "three_cordova"]
    # 	}]
    # 	cloud = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_cloud", "two_cloud", "three_cloud"]
    # 	}]
    # 	warehouse = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_warehouse", "two_warehouse", "three_warehouse"]
    # 	}]
    # 	shopify = [{
    # 		provider = "custom"
    # 		resolution_strategy = "and"
    # 		consents = ["one_shopify", "two_shopify", "three_shopify"]
    # 	}]
    # }
  }
}
