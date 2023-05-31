package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"gitlab.danalockapps.com/Backend/tools/terraform-provider-tanka.git/tanka"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tanka.Provider,
	})
}
