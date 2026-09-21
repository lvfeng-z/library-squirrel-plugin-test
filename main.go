package main

import (
	sdkdto "github.com/lvfeng-z/library-squirrel-sdk/dto"
	sdkplugin "github.com/lvfeng-z/library-squirrel-sdk/plugin"
)

func main() {
	handler := &TestTaskHandler{}
	fetcher := &TestSiteAuthorFetcher{}

	sdkplugin.Serve(sdkplugin.WithTaskHandler(handler),
		sdkplugin.WithSiteAuthorFetcher(fetcher),
		sdkplugin.WithActivate(func(ctx sdkdto.PluginContext) {
			Activate(ctx, handler, fetcher)
		}),
	)
}
