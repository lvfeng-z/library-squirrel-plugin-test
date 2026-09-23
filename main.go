package main

import (
	sdkdto "github.com/lvfeng-z/library-squirrel-sdk/dto"
	sdkplugin "github.com/lvfeng-z/library-squirrel-sdk/plugin"
)

func main() {
	handler := &TestWorkFetcher{}
	// 双拉取实例 fixture：与清单 siteAuthorFetch 两条目对应（同站点 bilibili、id 互异），
	// SDK 服务端按请求 extensionId 分派；应答带条目标记，实机可断言显选路由。
	mainFetcher := &TestSiteAuthorFetcher{entryID: "main"}
	altFetcher := &TestSiteAuthorFetcher{entryID: "alt"}

	sdkplugin.Serve(sdkplugin.WithWorkFetcher(handler),
		sdkplugin.WithSiteAuthorFetcher("main", mainFetcher),
		sdkplugin.WithSiteAuthorFetcher("alt", altFetcher),
		sdkplugin.WithActivate(func(ctx sdkdto.PluginContext) {
			Activate(ctx, handler, mainFetcher, altFetcher)
		}),
	)
}
