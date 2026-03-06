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

$ToolkitDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ToolkitDir
$SkillsSrc = Join-Path $ProjectRoot "skills"

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

Write-Host "📦 安装腾讯电子签合同智能 Skills"
Write-Host "   目标工具: $ToolName"
Write-Host "   安装目录: $DestDir"
Write-Host ""

foreach ($skill in $Skills) {
    $Src = Join-Path $SkillsSrc $skill
    if (-not (Test-Path $Src)) {
        Write-Host "⚠ 跳过 $skill（目录不存在）"
        continue
    }

    $SkillDest = Join-Path $DestDir $skill
    Write-Host "  → 安装 $skill..."
    New-Item -ItemType Directory -Force -Path $SkillDest | Out-Null

    # 复制 SKILL.md
    Copy-Item (Join-Path $Src "SKILL.md") -Destination $SkillDest -Force

    # 复制 scripts 目录
    $ScriptsDest = Join-Path $SkillDest "scripts"
    if (Test-Path $ScriptsDest) {
        Remove-Item -Recurse -Force $ScriptsDest
    }
    Copy-Item (Join-Path $Src "scripts") -Destination $SkillDest -Recurse -Force
}

Write-Host ""
Write-Host "✅ 安装完成！已安装 $($Skills.Count) 个 skills 到 $DestDir"
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
