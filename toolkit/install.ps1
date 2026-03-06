# 腾讯电子签合同智能 Skills — 一键安装脚本 (Windows PowerShell)
# 将 skills 安装到 AI 编码工具的 skills 目录
#
# 用法:
#   .\install.ps1 -Tool <claude|codebuddy|opencode> [-Target project|personal]
#
# 选项:
#   -Tool codebuddy   安装到 CodeBuddy skills 目录 (默认)
#   -Tool claude      安装到 Claude Code skills 目录
#   -Tool opencode    安装到 OpenCode skills 目录
#   -Target project   安装到当前项目目录 (默认)
#   -Target personal  安装到个人/全局目录
param(
    [ValidateSet("claude", "codebuddy", "opencode")]
    [string]$Tool = "codebuddy",

    [ValidateSet("project", "personal")]
    [string]$Target = "project"
)

$ErrorActionPreference = "Stop"

# 脚本所在目录（本地开发时在 toolkit/，打包后在 tarball 根目录）
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# 自动检测 skills 源目录
if (Test-Path (Join-Path $ScriptDir "contract-atoms")) {
    $SkillsSrc = $ScriptDir
} elseif (Test-Path (Join-Path (Split-Path -Parent $ScriptDir) "skills")) {
    $SkillsSrc = Join-Path (Split-Path -Parent $ScriptDir) "skills"
} else {
    Write-Error "找不到 skills 目录，请在安装包解压目录或项目根目录下运行"
    exit 1
}

# 根据工具和目标确定安装路径
switch ($Tool) {
    "claude" {
        $ToolDir = ".claude/skills"
        $PersonalDir = Join-Path $env:USERPROFILE ".claude/skills"
        $ToolName = "Claude Code"
    }
    "codebuddy" {
        $ToolDir = ".codebuddy/skills"
        $PersonalDir = Join-Path $env:USERPROFILE ".codebuddy/skills"
        $ToolName = "CodeBuddy"
    }
    "opencode" {
        $ToolDir = ".opencode/skills"
        $PersonalDir = Join-Path $env:USERPROFILE ".opencode/skills"
        $ToolName = "OpenCode"
    }
}

switch ($Target) {
    "project"  { $DestDir = Join-Path (Get-Location) $ToolDir }
    "personal" { $DestDir = $PersonalDir }
}

$Skills = @("contract-atoms", "contract-review", "contract-comparison")

# 检测当前平台
$Os = if ($IsLinux) { "linux" } elseif ($IsMacOS) { "darwin" } else { "windows" }
$Arch = if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -eq "Arm64") { "arm64" } else { "amd64" }
$Platform = "$Os-$Arch"

Write-Host "📦 安装腾讯电子签合同智能 Skills"
Write-Host "   目标工具: $ToolName"
Write-Host "   安装目录: $DestDir"
Write-Host "   当前平台: $Platform"
Write-Host ""

$Installed = 0
foreach ($skill in $Skills) {
    $Src = Join-Path $SkillsSrc $skill
    if (-not (Test-Path $Src)) {
        Write-Host "⚠ 跳过 $skill（目录不存在）"
        continue
    }

    $SkillDest = Join-Path $DestDir $skill
    Write-Host "  → 安装 $skill..."
    New-Item -ItemType Directory -Force -Path (Join-Path $SkillDest "scripts/bin") | Out-Null

    # 复制 SKILL.md
    Copy-Item (Join-Path $Src "SKILL.md") -Destination $SkillDest -Force

    # 复制 run 脚本
    Copy-Item (Join-Path $Src "scripts/run.sh") -Destination (Join-Path $SkillDest "scripts/") -Force
    $Ps1Path = Join-Path $Src "scripts/run.ps1"
    if (Test-Path $Ps1Path) {
        Copy-Item $Ps1Path -Destination (Join-Path $SkillDest "scripts/") -Force
    }

    # 复制二进制文件（区分平台子目录和平面结构）
    $BinCopied = 0
    $PlatformBinDir = Join-Path $Src "scripts/bin/$Platform"
    $FlatBinDir = Join-Path $Src "scripts/bin"

    if (Test-Path $PlatformBinDir) {
        # 本地开发：bin 按平台子目录组织
        Get-ChildItem -File $PlatformBinDir | ForEach-Object {
            Copy-Item $_.FullName -Destination (Join-Path $SkillDest "scripts/bin/") -Force
            $BinCopied++
        }
    } else {
        # tarball：bin 是平面结构
        Get-ChildItem -File $FlatBinDir -ErrorAction SilentlyContinue | ForEach-Object {
            Copy-Item $_.FullName -Destination (Join-Path $SkillDest "scripts/bin/") -Force
            $BinCopied++
        }
    }

    if ($BinCopied -eq 0) {
        Write-Host "    ⚠ 未找到 $Platform 的二进制文件（请先运行 make release）"
    } else {
        Write-Host "    ✅ 复制 $BinCopied 个二进制文件"
        $Installed++
    }
}

Write-Host ""
if ($Installed -eq 0) {
    Write-Host "⚠ 安装完成但缺少二进制文件。请先编译："
    Write-Host ""
    Write-Host "  cd $ScriptDir && make release"
    Write-Host ""
    Write-Host "然后重新运行 install.ps1"
} else {
    Write-Host "✅ 安装完成！已安装 $Installed/$($Skills.Count) 个 skills 到 $DestDir"
}
Write-Host ""
Write-Host "下一步：配置凭证（任选一种方式）："
Write-Host ""
Write-Host "  方式一（推荐）：创建 ~/.tsign/config.yaml"
Write-Host "    credentials:"
Write-Host "      secret_id: `"AKIDxxxxxxxx`""
Write-Host "      secret_key: `"xxxxxxxx`""
Write-Host "    operator:"
Write-Host "      user_id: `"yDwJxxx`""
Write-Host "    env: `"online`""
Write-Host ""
Write-Host "  方式二：设置环境变量（适合 CI/CD）"
Write-Host "    `$env:TENCENTCLOUD_SECRET_ID = `"AKIDxxxxxxxx`""
Write-Host "    `$env:TENCENTCLOUD_SECRET_KEY = `"xxxxxxxx`""
Write-Host "    `$env:ESS_OPERATOR_ID = `"yDwJxxx`""
