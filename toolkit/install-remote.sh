#!/usr/bin/env bash
# 腾讯电子签合同智能 Skills — 远程一键安装脚本
# 从 GitHub Releases 下载安装包并安装到 AI 编码工具的 skills 目录
#
# 快速安装:
#   curl -fsSL https://raw.githubusercontent.com/tencentess/ess-skills/main/toolkit/install-remote.sh | bash
#   curl -fsSL https://raw.githubusercontent.com/tencentess/ess-skills/main/toolkit/install-remote.sh | bash -s -- --tool=claude
#
# 用法:
#   bash install-remote.sh [选项]
#
# 选项:
#   --tool=<codebuddy|claude|opencode>  目标 AI 工具 (默认: codebuddy)
#   --target=<project|personal>         安装位置 (默认: project)
#   --version=<版本号>                   指定版本 (默认: latest)
#   --base-url=<URL>                    自定义下载地址（覆盖 GitHub Releases）
set -euo pipefail

# ============================================================
# 配置
# ============================================================
GITHUB_REPO="tencentess/ess-skills"
VERSION="latest"
TOOL="codebuddy"
TARGET="project"
BASE_URL=""

# ============================================================
# 解析参数
# ============================================================
for arg in "$@"; do
  case "$arg" in
    --tool=*)     TOOL="${arg#*=}" ;;
    --target=*)   TARGET="${arg#*=}" ;;
    --version=*)  VERSION="${arg#*=}" ;;
    --base-url=*) BASE_URL="${arg#*=}" ;;
    --help|-h)
      echo "腾讯电子签合同智能 Skills — 远程安装脚本"
      echo ""
      echo "快速安装:"
      echo "  curl -fsSL https://raw.githubusercontent.com/${GITHUB_REPO}/main/toolkit/install-remote.sh | bash"
      echo "  curl -fsSL https://raw.githubusercontent.com/${GITHUB_REPO}/main/toolkit/install-remote.sh | bash -s -- --tool=claude"
      echo ""
      echo "选项:"
      echo "  --tool=codebuddy   安装到 CodeBuddy skills 目录 (默认)"
      echo "  --tool=claude      安装到 Claude Code skills 目录"
      echo "  --tool=opencode    安装到 OpenCode skills 目录"
      echo "  --target=project   安装到当前项目目录 (默认)"
      echo "  --target=personal  安装到个人/全局目录"
      echo "  --version=latest   安装指定版本 (默认: latest)"
      echo "  --base-url=<URL>   自定义下载地址"
      exit 0
      ;;
    *)
      echo "❌ 未知参数: $arg" >&2
      echo "   运行 bash install-remote.sh --help 查看用法" >&2
      exit 1
      ;;
  esac
done

BASE_URL="${BASE_URL:-}"

# ============================================================
# 检测平台
# ============================================================
detect_platform() {
  local os arch
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"

  case "$os" in
    darwin)  os="darwin" ;;
    linux)   os="linux" ;;
    mingw*|msys*|cygwin*) os="windows" ;;
    *)
      echo "❌ 不支持的操作系统: $os" >&2
      exit 1
      ;;
  esac

  case "$arch" in
    x86_64|amd64)  arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *)
      echo "❌ 不支持的架构: $arch" >&2
      exit 1
      ;;
  esac

  echo "${os}-${arch}"
}

# ============================================================
# 检测下载工具
# ============================================================
detect_downloader() {
  if command -v curl &>/dev/null; then
    echo "curl"
  elif command -v wget &>/dev/null; then
    echo "wget"
  else
    echo "❌ 需要 curl 或 wget，请先安装" >&2
    exit 1
  fi
}

# ============================================================
# 下载文件
# ============================================================
download() {
  local url="$1" dest="$2" downloader
  downloader="$(detect_downloader)"

  echo "  ⬇ 下载: $url"
  case "$downloader" in
    curl) curl -fsSL --progress-bar -o "$dest" "$url" ;;
    wget) wget -q --show-progress -O "$dest" "$url" ;;
  esac
}

# 解析版本号（latest → 实际版本号）
resolve_version() {
  local ver="$1"
  if [ "$ver" = "latest" ]; then
    local api_url="https://api.github.com/repos/${GITHUB_REPO}/releases/latest"
    local downloader
    downloader="$(detect_downloader)"
    local tag
    case "$downloader" in
      curl) tag="$(curl -fsSL "$api_url" | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"//;s/".*//')" ;;
      wget) tag="$(wget -qO- "$api_url" | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"//;s/".*//')" ;;
    esac
    if [ -z "$tag" ]; then
      echo "❌ 无法获取最新版本号，请检查网络或指定 --version=<版本号>" >&2
      exit 1
    fi
    echo "${tag#v}"
  else
    echo "${ver#v}"
  fi
}

# ============================================================
# 主流程
# ============================================================
PLATFORM="$(detect_platform)"
echo "🔍 检测到平台: ${PLATFORM}"
echo ""

# 解析版本
RESOLVED_VERSION="$(resolve_version "$VERSION")"
echo "📌 版本: v${RESOLVED_VERSION}"

# 创建临时目录
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

# 构造下载 URL
ARCHIVE_NAME="tsign-skills-${RESOLVED_VERSION}-${PLATFORM}.tar.gz"
if [ -n "$BASE_URL" ]; then
  DOWNLOAD_URL="${BASE_URL}/${ARCHIVE_NAME}"
else
  DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/v${RESOLVED_VERSION}/${ARCHIVE_NAME}"
fi

echo "📦 下载腾讯电子签合同智能 Skills (v${RESOLVED_VERSION})"
download "$DOWNLOAD_URL" "${TMPDIR}/${ARCHIVE_NAME}"

# 解压
echo "📂 解压安装包..."
tar -xzf "${TMPDIR}/${ARCHIVE_NAME}" -C "$TMPDIR"

# 查找并执行 install.sh
INSTALL_SCRIPT="${TMPDIR}/install.sh"
if [ ! -f "$INSTALL_SCRIPT" ]; then
  # 可能在子目录中
  INSTALL_SCRIPT="$(find "$TMPDIR" -name 'install.sh' -maxdepth 2 | head -1)"
fi

if [ -z "$INSTALL_SCRIPT" ] || [ ! -f "$INSTALL_SCRIPT" ]; then
  echo "❌ 安装包中未找到 install.sh" >&2
  exit 1
fi

echo ""
bash "$INSTALL_SCRIPT" --tool="$TOOL" --target="$TARGET"
