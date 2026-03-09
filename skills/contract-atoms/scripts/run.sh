#!/usr/bin/env bash
# 合同智能原子操作 — 统一运行入口
# 用法: bash scripts/run.sh <command> [args...]
# 示例: bash scripts/run.sh review-create --resource-ids="xxx"
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SKILL_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# 检测平台和架构
detect_platform() {
  local os arch
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"
  case "$arch" in
    x86_64)  arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *) echo "❌ 不支持的架构: $arch" >&2; exit 1 ;;
  esac
  echo "${os}-${arch}"
}

if [ $# -lt 1 ]; then
  echo "用法: bash scripts/run.sh <command> [args...]" >&2
  echo "" >&2
  echo "可用命令:" >&2
  echo "  review-create      创建合同审查任务" >&2
  echo "  review-query       查询审查任务结果" >&2
  echo "  review-export      导出审查结果" >&2
  echo "  review-checklists  查看企业审查清单" >&2
  echo "  compare-create     创建合同对比任务" >&2
  echo "  compare-query      查询对比结果" >&2
  echo "  compare-export     导出对比报告" >&2
  exit 1
fi

CMD="$1"
shift

PLATFORM="$(detect_platform)"

# 查找可执行文件：优先安装后的扁平结构，再 fallback 到开发时的平台子目录
BIN="${SCRIPT_DIR}/bin/${CMD}"
if [ ! -f "$BIN" ]; then
  BIN="${SCRIPT_DIR}/bin/${PLATFORM}/${CMD}"
fi
if [ ! -f "$BIN" ]; then
  echo "❌ 未找到可执行文件: scripts/bin/${CMD} 或 scripts/bin/${PLATFORM}/${CMD}" >&2
  echo "   当前平台: ${PLATFORM}" >&2
  exit 1
fi

exec "$BIN" "$@"
