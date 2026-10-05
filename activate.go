package main

import (
	sdkdto "github.com/lvfeng-z/library-squirrel-sdk/dto"
)

// Activate 插件激活回调，在握手完成后由 SDK 调用。
// 扩展点注册与 URL 监听已由清单声明（宿主激活期派生），此处只做实例注入。
func Activate(ctx sdkdto.PluginContext, handler *TestWorkFetcher, fetchers ...*TestSiteAuthorFetcher) {
	handler.logger = ctx.GetLogger().Named("WorkFetch")
	for _, f := range fetchers {
		f.logger = ctx.GetLogger().Named("SiteAuthorFetcher(" + f.entryID + ")")
	}

	// 偏好域三方法演练回路（订阅前端演练页请求，经 SDK 偏好域执行并回发结果）
	RunPreferenceDrill(ctx)

	ctx.Infof("测试插件已激活")
}
