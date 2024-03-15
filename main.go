package main

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"

	"github.com/mach-composer/mach-composer-plugin-algolia/internal"
)

func main() {
	p := internal.NewAlgoliaPlugin()
	plugin.ServePlugin(p)
}
