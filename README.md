# library-squirrel-plugin-test

LibrarySquirrel 最小测试插件——用于验证前端扩展点（resourceViewer Handler 通道等），仅供开发测试，**不发布**。

## 用途

- 验证插件渲染器（resourceViewer）覆盖内置渲染器
- 验证插件↔主程序扩展点契约（props 注入、故障隔离等）
- 作为多扩展点测试的载体，按需在 `plugin.json` 的 `slots[]` 追加测试场景

## 当前测试场景

### test-article-viewer

声明 `resourceViewer`（`resourceType: "article"`），覆盖内置 ArticleRenderer。组件用 precompiled 工厂函数（注入 `__VUE__`），展示主程序注入的 `{resource, work}` 关键字段。

**验证**：主程序打开任意 article 资源 → 应显示"🔧 测试插件渲染器（article）" + resource/work 的 JSON 信息。卸载插件 → 回内置 ArticleRenderer。

## 构建

```powershell
pwsh ./build.ps1
```

打包 `plugin.json + views/` 到主程序 `../library-squirrel/resources/bundled-plugins/test-plugin.zip`。

主程序 `default_config.yaml` 已声明该 zip，`task dev` 启动时 `InstallBundledPlugins` 自动安装（首次安装后，DB 已记录；改组件后需在插件管理卸载重装或清 DB 才会重新解包）。

## 注意

- 禁止用 `import { X as Y }` 的 as 语法（Vite 工厂插件不兼容），需要别名直接在解构时改名。
- 样式遵循主题令牌契约（`var(--app-*)`），自动跟随主程序主题。
- 发布主程序前应移除：`test-plugin.zip` + `default_config.yaml` 中 test-plugin 声明 + 主程序内的 `test-plugin-src`（若残留）。
