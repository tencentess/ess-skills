#!/usr/bin/env bash
# 合同对比工作流 — 统一运行入口
# 用法: bash scripts/run.sh compare-workflow [args...]
# 示例: bash scripts/run.sh compare-workflow --origin-file="old.pdf" --diff-file="new.pdf"
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

CMD="${1:-compare-workflow}"
shift 2>/dev/null || true

PLATFORM="$(detect_platform)"
BIN="${SKILL_DIR}/bin/${PLATFORM}/${CMD}"

if [ ! -f "$BIN" ]; then
  echo "❌ 未找到可执行文件: bin/${PLATFORM}/${CMD}" >&2
  echo "   当前平台: ${PLATFORM}" >&2
  exit 1
fi

exec "$BIN" "$@"
