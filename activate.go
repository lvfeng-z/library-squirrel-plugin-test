package main

import (
	sdkdto "github.com/lvfeng-z/library-squirrel-sdk/dto"
)

// Activate 插件激活回调，在握手完成后由 SDK 调用。
// 扩展点注册与 URL 监听已由清单声明（宿主激活期派生），此处只做实例注入。
func Activate(ctx sdkdto.PluginContext, handler *TestTaskHandler, fetchers ...*TestSiteAuthorFetcher) {
	handler.logger = ctx.GetLogger().Named("TaskHandler")
	for _, f := range fetchers {
		f.logger = ctx.GetLogger().Named("SiteAuthorFetcher(" + f.entryID + ")")
	}

	ctx.Infof("测试插件已激活")
}
