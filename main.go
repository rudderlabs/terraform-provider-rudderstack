package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return rudderstack.New()
		},
	}

	if debugMode {
		err := plugin.Debug(context.Background(), "rudderlabs/rudderstack", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
