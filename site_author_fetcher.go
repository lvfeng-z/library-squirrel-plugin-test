package main

import (
	"context"
	"fmt"
	"strings"

	sdkdto "github.com/lvfeng-z/library-squirrel-sdk/dto"
)

// TestSiteAuthorFetcher 站点作者信息拉取桩（双实例 fixture）。只回一块 meta，不产头像字节块
// （无头像是契约允许的合法态）；每个实例携带自己的条目 id，应答作者名带「[条目id]」标记，
// 供实机验收断言「冲突选择器显选后请求路由到被选实例」——两实例应答可区分。
type TestSiteAuthorFetcher struct {
	entryID string // 对应清单 extensions.siteAuthorFetch 条目 id（main/alt）
	logger  sdkdto.Logger
}

// FetchSiteAuthorInfo 按站点身份键回作者元数据。归属由清单 extensions.siteAuthorFetch.sites
// 声明，主程序按站点收窄候选并按条目 id 分派（请求携带 extensionId），本实现只处理被路由到的请求。
func (f *TestSiteAuthorFetcher) FetchSiteAuthorInfo(ctx context.Context, req *sdkdto.FetchSiteAuthorInfoRequest, send func(chunk *sdkdto.AuthorInfoChunk) error) error {
	siteAuthorId := strings.TrimSpace(req.GetSiteAuthorId())
	if got := req.GetExtensionId(); got != f.entryID {
		// 宿主分派错误才会走到这里（正常时 SDK 服务端已按 extensionId 选中本实例）
		f.logger.Warnf("路由条目不匹配：本实例=%s 请求extensionId=%s", f.entryID, got)
	}
	f.logger.Infof("测试插件应答站点作者拉取请求 entryID=%s siteAuthorId=%s", f.entryID, siteAuthorId)
	return send(&sdkdto.AuthorInfoChunk{Payload: &sdkdto.AuthorInfoChunk_Meta{
		Meta: &sdkdto.AuthorInfoMeta{
			AuthorName: fmt.Sprintf("测试插件作者[%s] %s", f.entryID, siteAuthorId),
			Introduce:  fmt.Sprintf("测试插件双实例作者信息桩（条目 %s）", f.entryID),
			Homepage:   "https://space.bilibili.com/" + siteAuthorId,
		},
	}})
}
