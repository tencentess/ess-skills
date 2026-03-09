# 合同对比工作流 — 统一运行入口 (Windows PowerShell)
# 用法: .\scripts\run.ps1 compare-workflow [args...]
# 示例: .\scripts\run.ps1 compare-workflow --origin-file="old.pdf" --diff-file="new.pdf"
param(
    [Parameter(Position=0)]
    [string]$Command = "compare-workflow",
    [Parameter(Position=1, ValueFromRemainingArguments=$true)]
    [string[]]$Arguments
)

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$SkillDir = Split-Path -Parent $ScriptDir

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
