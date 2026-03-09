# 腾讯电子签合同智能 Skills — 远程一键安装脚本 (Windows PowerShell)
# 从 GitHub Releases 下载安装包并安装到 AI 编码工具的 skills 目录
#
# 快速安装:
#   irm https://raw.githubusercontent.com/tencentess/ess-skills/main/toolkit/install-remote.ps1 -OutFile install-remote.ps1; powershell -ExecutionPolicy Bypass -File .\install-remote.ps1
#
# 用法:
#   .\install-remote.ps1 [-Tool codebuddy] [-Target project] [-Version latest] [-BaseUrl <URL>]
param(
    [ValidateSet("claude", "codebuddy", "opencode")]
    [string]$Tool = "codebuddy",

    [ValidateSet("project", "personal")]
    [string]$Target = "project",

    [string]$Version = "latest",

    [string]$BaseUrl = ""
)

$ErrorActionPreference = "Stop"

$GithubRepo = "tencentess/ess-skills"

# ============================================================
# 检测平台
# ============================================================
$Arch = switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture) {
    "X64"   { "amd64" }
    "Arm64" { "arm64" }
    default { Write-Error "不支持的架构: $_"; exit 1 }
}

$Platform = "windows-$Arch"
Write-Host "🔍 检测到平台: $Platform"
Write-Host ""

# ============================================================
# 下载安装包
# ============================================================
# 解析 latest 版本号
if ($Version -eq "latest") {
    $ApiUrl = "https://api.github.com/repos/$GithubRepo/releases/latest"
    $Release = Invoke-RestMethod -Uri $ApiUrl -UseBasicParsing
    $ResolvedVersion = $Release.tag_name -replace '^v', ''
    if (-not $ResolvedVersion) {
        Write-Error "无法获取最新版本号，请检查网络或指定 -Version <版本号>"
        exit 1
    }
} else {
    $ResolvedVersion = $Version -replace '^v', ''
}

Write-Host "📌 版本: v$ResolvedVersion"

$ArchiveName = "tsign-skills-$ResolvedVersion-$Platform.tar.gz"
if ($BaseUrl) {
    $DownloadUrl = "$BaseUrl/$ArchiveName"
} else {
    $DownloadUrl = "https://github.com/$GithubRepo/releases/download/v$ResolvedVersion/$ArchiveName"
}

$TmpDir = Join-Path ([System.IO.Path]::GetTempPath()) "tsign-install-$(Get-Random)"
New-Item -ItemType Directory -Force -Path $TmpDir | Out-Null

try {
    $ArchivePath = Join-Path $TmpDir $ArchiveName

    Write-Host "📦 下载腾讯电子签合同智能 Skills (v$ResolvedVersion)"
    Write-Host "  ⬇ 下载: $DownloadUrl"
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $ArchivePath -UseBasicParsing

    # 解压
    Write-Host "📂 解压安装包..."
    tar -xzf $ArchivePath -C $TmpDir

    # 查找并执行 install.ps1
    $InstallScript = Get-ChildItem -Path $TmpDir -Filter "install.ps1" -Recurse -Depth 2 | Select-Object -First 1

    if (-not $InstallScript) {
        # 回退到 install.sh（如果用户有 bash 可用）
        Write-Error "安装包中未找到 install.ps1"
        exit 1
    }

    Write-Host ""
    & powershell.exe -ExecutionPolicy Bypass -File $InstallScript.FullName -Tool $Tool -Target $Target

} finally {
    # 清理临时目录
    if (Test-Path $TmpDir) {
        Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
    }
}
