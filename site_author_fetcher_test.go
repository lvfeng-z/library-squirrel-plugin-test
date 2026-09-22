package main

import (
	"context"
	"strings"
	"testing"

	sdkdto "github.com/lvfeng-z/library-squirrel-sdk/dto"
)

// fetchLogStub 记录 Warnf 调用，用于断言路由不匹配告警路径
type fetchLogStub struct {
	warns []string
}

func (l *fetchLogStub) Infof(template string, args ...any)  {}
func (l *fetchLogStub) Debugf(template string, args ...any) {}
func (l *fetchLogStub) Errorf(template string, args ...any) {}
func (l *fetchLogStub) Warnf(template string, args ...any) {
	l.warns = append(l.warns, template)
}
func (l *fetchLogStub) Named(name string) sdkdto.Logger { return l }

// TestDualFetcherInstancesDistinguishable 双实例 fixture 锚定：main/alt 两实现应答携带
// 各自条目标记（「[main]」/「[alt]」），实机显选后可据应答标记断言路由到被选实例。
// 对应清单 plugin.json extensions.siteAuthorFetch 两条目（id=main/alt，同站点 bilibili）。
func TestDualFetcherInstancesDistinguishable(t *testing.T) {
	log := &fetchLogStub{}
	instances := map[string]*TestSiteAuthorFetcher{
		"main": {entryID: "main", logger: log},
		"alt":  {entryID: "alt", logger: log},
	}

	names := make(map[string]string)
	for entryID, fetcher := range instances {
		var chunks []*sdkdto.AuthorInfoChunk
		send := func(chunk *sdkdto.AuthorInfoChunk) error { chunks = append(chunks, chunk); return nil }
		req := &sdkdto.FetchSiteAuthorInfoRequest{
			SiteKey:      "bilibili",
			SiteAuthorId: "2233",
			ExtensionId:  entryID,
		}
		if err := fetcher.FetchSiteAuthorInfo(context.Background(), req, send); err != nil {
			t.Fatalf("entryID=%s 拉取失败: %v", entryID, err)
		}
		if len(chunks) != 1 || chunks[0].GetMeta() == nil {
			t.Fatalf("entryID=%s 应只回一块 meta", entryID)
		}
		names[entryID] = chunks[0].GetMeta().GetAuthorName()
		wantMarker := "[" + entryID + "]"
		if !strings.Contains(names[entryID], wantMarker) {
			t.Fatalf("entryID=%s 应答作者名 %q 应含标记 %q", entryID, names[entryID], wantMarker)
		}
	}
	if names["main"] == names["alt"] {
		t.Fatalf("两实例应答应可区分，实际相同：%q", names["main"])
	}
	if len(log.warns) != 0 {
		t.Fatalf("条目 id 匹配时不应告警路由不匹配，实际告警 %d 次", len(log.warns))
	}
}

// TestFetcherWarnsOnMismatchedRoute 路由不匹配（宿主分派错误）时告警但不失败
func TestFetcherWarnsOnMismatchedRoute(t *testing.T) {
	log := &fetchLogStub{}
	fetcher := &TestSiteAuthorFetcher{entryID: "main", logger: log}
	send := func(chunk *sdkdto.AuthorInfoChunk) error { return nil }
	req := &sdkdto.FetchSiteAuthorInfoRequest{SiteAuthorId: "1", ExtensionId: "alt"}
	if err := fetcher.FetchSiteAuthorInfo(context.Background(), req, send); err != nil {
		t.Fatalf("不匹配路由不应报错: %v", err)
	}
	if len(log.warns) != 1 {
		t.Fatalf("应告警一次路由不匹配，实际 %d 次", len(log.warns))
	}
}
