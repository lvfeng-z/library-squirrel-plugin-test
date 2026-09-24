// 测试插件 resolver：『全量设置 → 条目参与度』纯函数（宿主内嵌 JS 引擎执行，
// 契约类型见 SDK settingresolver/contract.d.ts）。布尔开关 enableParticipation 控制
// resourceViewer 参与（frontendExtensions 条目 test-article-viewer）与作品拉取候选
// （workFetch 条目 main，URL 监听随条目联动摘除/恢复）。
//
// 语法面锚定 ES2015（const/箭头/模板串/解构可用；async/await、可选链等 ES2017+ 勿用）。
// 输入 settings 值以存储形态喂入（boolean 设置为字符串 "true"/"false"），比较前用
// String() 归一；缺键/非 "false" 取值 = 基线参与（与声明默认值 "true" 一致）。
function resolve(input) {
	var settings = input && input.settings ? input.settings : {};
	var participate = String(settings.enableParticipation) !== "false";
	if (participate) {
		return { version: 1, entries: [] }; // 快照语义：未列出条目 = 基线参与
	}
	var reason = "设置「参与派生面」已关闭";
	return {
		version: 1,
		entries: [
			{ point: "frontendExtensions", id: "test-article-viewer", active: false, reason: reason },
			{ point: "workFetch", id: "main", active: false, reason: reason }
		]
	};
}
