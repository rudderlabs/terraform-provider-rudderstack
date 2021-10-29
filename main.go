package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-provider-rudderstack/rudderstack"
)

func main() {
	log.Println("Terraform Provider Rudderstack started.")

	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &tfsdk.ServeOpts{
		Name: "rudderstack",
	}

	/*
	if debugMode {
		err := tfsdk.Debug(context.Background(), rudderstack.New, opts)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	*/

	tfsdk.Serve(context.Background(), rudderstack.New, *opts)
}
