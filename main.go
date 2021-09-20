package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-provider-rudderstack/rudderstack"
)


func main() {
	tfsdk.Serve(context.Background(), rudderstack.New, tfsdk.ServeOpts{
		Name: "rudderstack",
	})
}
