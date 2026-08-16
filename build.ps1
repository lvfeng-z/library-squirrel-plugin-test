# 打包测试插件 zip 到主程序 resources/bundled-plugins/
# 用法：仓库根目录运行  pwsh ./build.ps1  或  powershell ./build.ps1
$ErrorActionPreference = 'Stop'
$repoRoot = $PSScriptRoot
# Main repo is the sibling directory (Split-Path normalizes; a literal "..\" makes Compress-Archive silently fail)
$mainRepo = Join-Path (Split-Path $repoRoot -Parent) 'library-squirrel'
$dest = Join-Path $mainRepo 'resources\bundled-plugins\test-plugin.zip'

# Build identity stamping: inject git describe (source-state id) into the packaged plugin.json as "buildId".
# This repo has no dist staging: copy plugin.json to a temp dir, stamp, then zip; repo source plugin.json stays untouched.
# Keep this block ASCII-only: powershell -File decodes no-BOM scripts as GBK, and CJK comments can swallow the next line.
# Text insertion instead of ConvertTo-Json round-trip (avoids PS5.1 depth/escape mangling); write UTF-8 without BOM (Go json rejects BOM).
$buildId = git describe --tags --always --dirty
if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($buildId)) {
    Write-Host "Build failed! git describe failed (not a git repo?), buildId is required." -ForegroundColor Red
    exit 1
}
$stage = Join-Path $env:TEMP ("test-plugin-stage-" + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $stage | Out-Null
try {
    Copy-Item (Join-Path $repoRoot 'plugin.json') $stage
    $manifestPath = Join-Path $stage 'plugin.json'
    $manifestText = [System.IO.File]::ReadAllText($manifestPath)
    $stamped = $manifestText -replace '^\s*\{', ('{"buildId": "' + $buildId.Trim() + '",')
    [System.IO.File]::WriteAllText($manifestPath, $stamped, (New-Object System.Text.UTF8Encoding($false)))
    Write-Host "  buildId: $($buildId.Trim())"

    # Pack plugin.json + views to zip root (same layout as other bundled zips: plugin.json at root, views flattened).
    # Suppress progress output before Compress-Archive (same workaround as pixiv/bilibili build scripts):
    # no-BOM CJK scripts decoded as GBK by powershell -File make Write-Progress break Compress-Archive silently.
    $ProgressPreference = 'SilentlyContinue'
    Compress-Archive -Path $manifestPath, (Join-Path $repoRoot 'views') -DestinationPath $dest -Force
} finally {
    Remove-Item -Recurse -Force $stage -ErrorAction SilentlyContinue
}
Write-Host "Packed test-plugin to: $dest"
