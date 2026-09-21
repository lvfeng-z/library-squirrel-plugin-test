package main

import (
	"context"
	"strings"

	sdkdto "github.com/lvfeng-z/library-squirrel-sdk/dto"
)

// TestSiteAuthorFetcher 站点作者信息拉取桩。只回一块 meta，不产头像字节块
// （无头像是契约允许的合法态）；作者名带插件标记，便于分辨应答来源。
type TestSiteAuthorFetcher struct {
	logger sdkdto.Logger
}

// FetchSiteAuthorInfo 按站点身份键回作者元数据。归属由清单 extensions.siteAuthorFetch.sites
// 声明，主程序按站点收窄候选，本实现只处理被路由到的请求。
func (f *TestSiteAuthorFetcher) FetchSiteAuthorInfo(ctx context.Context, req *sdkdto.FetchSiteAuthorInfoRequest, send func(chunk *sdkdto.AuthorInfoChunk) error) error {
	siteAuthorId := strings.TrimSpace(req.GetSiteAuthorId())
	f.logger.Infof("测试插件应答站点作者拉取请求 siteAuthorId=%s", siteAuthorId)
	return send(&sdkdto.AuthorInfoChunk{Payload: &sdkdto.AuthorInfoChunk_Meta{
		Meta: &sdkdto.AuthorInfoMeta{
			AuthorName: "测试插件作者 " + siteAuthorId,
			Introduce:  "测试插件返回的作者信息桩",
			Homepage:   "https://space.bilibili.com/" + siteAuthorId,
		},
	}})
}
