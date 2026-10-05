package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	sdkdto "github.com/lvfeng-z/library-squirrel-sdk/dto"
)

// 偏好三方法演练回路：前端演练页（view 扩展 test-preference-drill）经 Wails 事件发请求，
// 本回路在插件进程内经 SDK PluginContext 偏好域三方法（GetPreference/SetPreference/
// ListMyPreferences）执行并把结果回发前端展示——验证宿主侧偏好域实现的真实链路。
// 通道形态见 doc/plugin-dev-guide.md 第九章：前端 Events.Emit → 宿主桥 → 插件订阅；
// 插件 PublishToFrontend → 宿主桥 → 前端 Events.On。
const (
	preferenceDrillReqTopic  = "plugin:test-plugin:preference:request"
	preferenceDrillRespTopic = "plugin:test-plugin:preference:response"
	// preferenceDrillKey 演练固定键：写入/覆写共用同键，演示整值覆写语义
	preferenceDrillKey = "download.form"
)

// preferenceDrillRequest 演练请求（前端按钮发出）。action 取值：set（写入）/ overwrite
// （同键覆写）/ get（读取）/ list（列键）；seq 为前端自增序号，回显用于结果配对。
type preferenceDrillRequest struct {
	Action string `json:"action"`
	Seq    int64  `json:"seq"`
}

// preferenceDrillResponse 演练结果（回发前端展示）。get 动作以 found 区分有记录/无记录；
// list 动作携带键清单；失败时 ok=false 且 error 带原因。
type preferenceDrillResponse struct {
	Action      string   `json:"action"`
	Seq         int64    `json:"seq"`
	OK          bool     `json:"ok"`
	Error       string   `json:"error,omitempty"`
	Found       bool     `json:"found"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Data        string   `json:"data,omitempty"`
	Keys        []string `json:"keys"`
}

// RunPreferenceDrill 启动偏好演练回路：订阅演练请求 topic，逐条经 SDK 偏好域执行并回发
// 结果。激活期调用一次，随插件进程存续；订阅失败仅记日志不阻断激活（演练页会超时提示）。
func RunPreferenceDrill(ctx sdkdto.PluginContext) {
	ch, err := ctx.SubscribeFrontend(preferenceDrillReqTopic)
	if err != nil {
		ctx.Errorf("偏好演练订阅失败: %v", err)
		return
	}
	ctx.Infof("偏好演练[就绪] topic=%s key=%s", preferenceDrillReqTopic, preferenceDrillKey)
	go func() {
		for data := range ch {
			resp := handlePreferenceDrill(ctx, data)
			raw, err := json.Marshal(resp)
			if err != nil {
				ctx.Errorf("偏好演练结果序列化失败: %v", err)
				continue
			}
			if err := ctx.PublishToFrontend(preferenceDrillRespTopic, raw); err != nil {
				ctx.Errorf("偏好演练结果回发失败: %v", err)
			}
		}
	}()
}

// handlePreferenceDrill 单条演练请求的执行：按动作分派到 SDK 偏好域三方法，
// 每步落一条可 grep 的演练日志（server.log 内 Plugin[测试插件] 作用域）。
func handlePreferenceDrill(ctx sdkdto.PluginContext, data []byte) preferenceDrillResponse {
	var req preferenceDrillRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return preferenceDrillResponse{OK: false, Error: "请求解析失败: " + err.Error(), Keys: []string{}}
	}
	resp := preferenceDrillResponse{Action: req.Action, Seq: req.Seq, Keys: []string{}}
	switch req.Action {
	case "set", "overwrite":
		value := drillValue(req.Action, req.Seq)
		if err := ctx.SetPreference(preferenceDrillKey, value); err != nil {
			resp.Error = fmt.Sprintf("SetPreference 失败: %v", err)
			ctx.Errorf("偏好演练[%s] SetPreference 失败: %v", drillActionLabel(req.Action), err)
			return resp
		}
		ctx.Infof("偏好演练[%s] SetPreference key=%s title=%s data=%s",
			drillActionLabel(req.Action), preferenceDrillKey, value.Title, value.Data)
		resp.OK = true
		resp.Title = value.Title
		resp.Data = value.Data
	case "get":
		v, found, err := ctx.GetPreference(preferenceDrillKey)
		if err != nil {
			resp.Error = fmt.Sprintf("GetPreference 失败: %v", err)
			ctx.Errorf("偏好演练[读取] GetPreference 失败: %v", err)
			return resp
		}
		resp.OK = true
		resp.Found = found
		if found {
			resp.Title = v.Title
			resp.Description = v.Description
			resp.Data = v.Data
			ctx.Infof("偏好演练[读取] key=%s found=true title=%s data=%s", preferenceDrillKey, v.Title, v.Data)
		} else {
			ctx.Infof("偏好演练[读取] key=%s found=false（无记录，插件下次将重新发起问答）", preferenceDrillKey)
		}
	case "list":
		keys, err := ctx.ListMyPreferences()
		if err != nil {
			resp.Error = fmt.Sprintf("ListMyPreferences 失败: %v", err)
			ctx.Errorf("偏好演练[列键] ListMyPreferences 失败: %v", err)
			return resp
		}
		if keys == nil {
			keys = []string{}
		}
		resp.OK = true
		resp.Keys = keys
		ctx.Infof("偏好演练[列键] keys=%v", keys)
	default:
		resp.Error = "未知动作: " + req.Action
	}
	return resp
}

// drillValue 构造演练写入值：写入与覆写同键不同 title/data（data 携带 seq 供区分轮次），
// 值信封按 SDK 约定四字段（schemaVersion 初值 1）。
func drillValue(action string, seq int64) *sdkdto.PreferenceValue {
	if action == "overwrite" {
		return &sdkdto.PreferenceValue{
			SchemaVersion: 1,
			Title:         "测试偏好：图文形态·文章模式",
			Description:   "测试插件偏好演练覆写",
			Data:          `{"form":"article","drillSeq":` + strconv.FormatInt(seq, 10) + `}`,
		}
	}
	return &sdkdto.PreferenceValue{
		SchemaVersion: 1,
		Title:         "测试偏好：图文形态",
		Description:   "测试插件偏好演练写入",
		Data:          `{"form":"multiImage","drillSeq":` + strconv.FormatInt(seq, 10) + `}`,
	}
}

// drillActionLabel 演练日志的动作中文标签
func drillActionLabel(action string) string {
	if action == "overwrite" {
		return "覆写"
	}
	return "写入"
}
