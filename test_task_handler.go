package main

import (
	"context"
	"errors"
	"strings"

	sdkdto "github.com/lvfeng-z/library-squirrel-sdk/dto"
)

// errStubNoDownload 任务执行面的统一应答：本插件只做声明与路由应答，不产生任何资源。
var errStubNoDownload = errors.New("测试插件不执行下载，请改选 bilibiliSuite 处理该任务")

// TestTaskHandler 任务处理器桩。Create/CreateWorkInfo 返回合法的最小响应以支撑任务创建与
// 冲突候选选择；Start/Retry/Resume 一律拒绝，避免桩插件产出残缺资源。
type TestTaskHandler struct {
	logger sdkdto.Logger
}

// Create 由 URL 派生站点侧作品 ID，返回单条任务创建响应（站点键 bilibili、资源类型 video）。
func (h *TestTaskHandler) Create(url string) (*sdkdto.TaskCreateResult, error) {
	siteWorkId := deriveSiteWorkId(url)
	h.logger.Infof("测试插件创建任务 siteWorkId=%s url=%s", siteWorkId, url)
	return sdkdto.BatchResult([]*sdkdto.TaskCreateResponse{{
		PluginTaskId: "test-" + siteWorkId,
		TaskName:     "测试插件任务 " + siteWorkId,
		SiteWorkId:   siteWorkId,
		Url:          url,
		SiteName:     "bilibili",
		SiteKey:      "bilibili",
		ResourceType: sdkdto.ResourceTypeVideo,
		InvolvedRoles: []string{
			sdkdto.StoreRoleVideoTrack, sdkdto.StoreRoleAudioTrack, sdkdto.StoreRoleThumbnail,
		},
	}}), nil
}

// CreateWorkInfo 回填作品信息：作品身份取自任务行的站点与站点侧作品 ID。
func (h *TestTaskHandler) CreateWorkInfo(task *sdkdto.TaskDTO) (*sdkdto.WorkResponse, error) {
	return testWorkResponse(task), nil
}

// Start 拒绝执行：桩不产出资源。
func (h *TestTaskHandler) Start(ctx context.Context, task *sdkdto.TaskDTO, storeRoles []string) ([]*sdkdto.StoreSpec, *sdkdto.WorkResponse, error) {
	return nil, nil, errStubNoDownload
}

// Retry 拒绝执行：桩不产出资源。
func (h *TestTaskHandler) Retry(task *sdkdto.TaskDTO) (*sdkdto.WorkResponse, error) {
	return nil, errStubNoDownload
}

// Pause 无在途执行，空操作成功。
func (h *TestTaskHandler) Pause(param *sdkdto.TaskResParam) error {
	return nil
}

// Stop 无在途执行，空操作成功。
func (h *TestTaskHandler) Stop(param *sdkdto.TaskResParam) error {
	return nil
}

// Resume 拒绝续传：桩不产出资源，无可续传轨道。
func (h *TestTaskHandler) Resume(ctx context.Context, param *sdkdto.TaskResumeParam) ([]*sdkdto.StoreSpec, *sdkdto.WorkResponse, error) {
	return nil, nil, errStubNoDownload
}

// testWorkResponse 构造最小作品信息响应。
func testWorkResponse(task *sdkdto.TaskDTO) *sdkdto.WorkResponse {
	siteWorkId := "unknown"
	if task != nil && task.SiteWorkId != nil && *task.SiteWorkId != "" {
		siteWorkId = *task.SiteWorkId
	}
	return &sdkdto.WorkResponse{
		Work: &sdkdto.WorkDTO{
			SiteWorkId:   &siteWorkId,
			SiteWorkName: ptrOf("测试插件作品 " + siteWorkId),
		},
		Site: &sdkdto.SiteDTO{
			SiteKey:  "bilibili",
			SiteName: ptrOf("bilibili"),
		},
		SiteAuthors: []*sdkdto.TaskSiteAuthorDTO{{
			SiteAuthorId: siteWorkId,
			AuthorName:   "测试插件作者",
			SiteKey:      "bilibili",
		}},
	}
}

// ptrOf 取字符串地址（proto 消息的可空标量字段以指针表达）
func ptrOf(s string) *string {
	return &s
}

// deriveSiteWorkId 取 URL 最后一个路径段作为站点侧作品 ID（只做标识来源，不做站点解析）。
func deriveSiteWorkId(url string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(url), "/")
	if i := strings.LastIndexByte(trimmed, '/'); i >= 0 && i < len(trimmed)-1 {
		return trimmed[i+1:]
	}
	return "unknown"
}
