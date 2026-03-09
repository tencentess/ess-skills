# 合同智能原子操作 — 统一运行入口 (Windows PowerShell)
# 用法: .\scripts\run.ps1 <command> [args...]
# 示例: .\scripts\run.ps1 review-create --resource-ids="xxx"
param(
    [Parameter(Position=0)]
    [string]$Command,
    [Parameter(Position=1, ValueFromRemainingArguments=$true)]
    [string[]]$Arguments
)

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$SkillDir = Split-Path -Parent $ScriptDir

if (-not $Command) {
    Write-Error @"
用法: .\scripts\run.ps1 <command> [args...]

可用命令:
  review-create      创建合同审查任务
  review-query       查询审查任务结果
  review-export      导出审查结果
  review-checklists  查看企业审查清单
  compare-create     创建合同对比任务
  compare-query      查询对比结果
  compare-export     导出对比报告
"@
    exit 1
}

# 检测架构
$Arch = switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture) {
    "X64"   { "amd64" }
    "Arm64" { "arm64" }
    default { Write-Error "不支持的架构: $_"; exit 1 }
}

$Platform = "windows-$Arch"

# 查找可执行文件：优先安装后的扁平结构，再 fallback 到开发时的平台子目录
$Bin = $null
foreach ($candidate in @(
    (Join-Path $ScriptDir "bin" "$Command.exe"),
    (Join-Path $ScriptDir "bin" $Command),
    (Join-Path $ScriptDir "bin" $Platform "$Command.exe"),
    (Join-Path $ScriptDir "bin" $Platform $Command)
)) {
    if (Test-Path $candidate) { $Bin = $candidate; break }
}
if (-not $Bin) {
    Write-Error "未找到可执行文件: scripts/bin/$Command.exe 或 scripts/bin/$Platform/$Command.exe`n当前平台: $Platform"
    exit 1
}

& $Bin @Arguments
