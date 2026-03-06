# 合同审查工作流 — 统一运行入口 (Windows PowerShell)
# 用法: .\scripts\run.ps1 review-workflow [args...]
# 示例: .\scripts\run.ps1 review-workflow --file="/path/to/contract.pdf" --policy-type=0
param(
    [Parameter(Position=0)]
    [string]$Command = "review-workflow",
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
$Bin = Join-Path $SkillDir "bin" $Platform "$Command.exe"

if (-not (Test-Path $Bin)) {
    $BinNoExt = Join-Path $SkillDir "bin" $Platform $Command
    if (Test-Path $BinNoExt) {
        $Bin = $BinNoExt
    } else {
        Write-Error "未找到可执行文件: bin/$Platform/$Command.exe`n当前平台: $Platform"
        exit 1
    }
}

& $Bin @Arguments
