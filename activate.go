package main

import (
	sdkdto "github.com/lvfeng-z/library-squirrel-sdk/dto"
)

// Activate 插件激活回调，在握手完成后由 SDK 调用。注册任务处理器与 URL 监听器。
func Activate(ctx sdkdto.PluginContext, handler *TestTaskHandler, fetcher *TestSiteAuthorFetcher) {
	handler.logger = ctx.GetLogger().Named("TaskHandler")
	fetcher.logger = ctx.GetLogger().Named("SiteAuthorFetcher")

	if err := ctx.RegisterTaskHandler("main", "TestTaskHandler",
		"测试任务处理器（只声明与应答，不做真实抓取）", handler); err != nil {
		ctx.Errorf("注册 TaskHandler 失败: %v", err)
		return
	}

	if err := ctx.RegisterUrlListener("main", urlListenerPatterns()); err != nil {
		ctx.Errorf("注册 URL 监听器失败: %v", err)
		return
	}

	ctx.Infof("测试插件已激活")
}

// urlListenerPatterns 任务路由的 URL 监听正则。与 bilibili 套件的模式同域同形，
// 令同一 bilibili URL 同时命中两个插件，任务创建进入候选选择面。
func urlListenerPatterns() []string {
	return []string{
		`^https?://www\.bilibili\.com/video/(BV[a-zA-Z0-9]+|av\d+)`,
		`^https?://b23\.tv/[A-Za-z0-9]+`,
		`^https?://www\.bilibili\.com/(opus|dynamic)/\d+`,
		`^https?://www\.bilibili\.com/read/cv\d+`,
	}
}
