# 打包测试插件 zip 到主程序 resources/bundled-plugins/
# 用法：仓库根目录运行  pwsh ./build.ps1  或  powershell ./build.ps1
$ErrorActionPreference = 'Stop'
$repoRoot = $PSScriptRoot
# Main repo is the sibling directory (Split-Path normalizes; a literal "..\" makes Compress-Archive silently fail)
$mainRepo = Join-Path (Split-Path $repoRoot -Parent) 'library-squirrel'
$dest = Join-Path $mainRepo 'resources\bundled-plugins\test-plugin.zip'
# 打包 plugin.json + views 到 zip 根（与现有 bundled zip 结构一致：plugin.json 在根、views 平铺）
Compress-Archive -Path (Join-Path $repoRoot 'plugin.json'), (Join-Path $repoRoot 'views') -DestinationPath $dest -Force
Write-Host "Packed test-plugin to: $dest"
