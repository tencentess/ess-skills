#!/usr/bin/env bash
# 腾讯电子签合同智能 Skills — 发布打包脚本
# 将 skills/ 目录按平台拆分，打包为 tsign-skills-{version}-{platform}.tar.gz
#
# 用法:
#   bash scripts/package.sh [--version=<版本号>] [--output=<输出目录>]
#
# 前置条件:
#   已运行 make release（skills/ 目录下已有各 skill 的 SKILL.md + scripts/）
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TOOLKIT_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$TOOLKIT_DIR")"
SKILLS_DIR="${PROJECT_ROOT}/skills"
VERSION="1.0.0"
OUTPUT_DIR="${PROJECT_ROOT}/dist"

# ============================================================
# 解析参数
# ============================================================
for arg in "$@"; do
  case "$arg" in
    --version=*)  VERSION="${arg#*=}" ;;
    --output=*)   OUTPUT_DIR="${arg#*=}" ;;
    --help|-h)
      echo "用法: bash scripts/package.sh [--version=<版本号>] [--output=<输出目录>]"
      echo ""
      echo "选项:"
      echo "  --version=1.0.0   版本号 (默认: 1.0.0)"
      echo "  --output=dist     输出目录 (默认: dist)"
      echo ""
      echo "前置条件: 先运行 make release"
      exit 0
      ;;
    *)
      echo "❌ 未知参数: $arg" >&2
      exit 1
      ;;
  esac
done

SKILLS=(contract-atoms contract-review contract-comparison)

# 支持的平台列表（与 Go 交叉编译目标一致）
PLATFORMS=(
  "darwin-amd64"
  "darwin-arm64"
  "linux-amd64"
  "linux-arm64"
  "windows-amd64"
  "windows-arm64"
)

# ============================================================
# 检查 skills 目录
# ============================================================
for skill in "${SKILLS[@]}"; do
  if [ ! -d "${SKILLS_DIR}/${skill}" ]; then
    echo "❌ 缺少 ${SKILLS_DIR}/${skill}，请先运行 make release" >&2
    exit 1
  fi
done

# ============================================================
# 准备输出目录
# ============================================================
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

echo "📦 打包腾讯电子签合同智能 Skills v${VERSION}"
echo "   输出目录: ${OUTPUT_DIR}"
echo ""

# ============================================================
# 辅助函数：获取平台对应的二进制后缀
# ============================================================
bin_suffix() {
  local platform="$1"
  case "$platform" in
    windows-*) echo ".exe" ;;
    *)         echo "" ;;
  esac
}

# 将 platform 格式 (darwin-amd64) 转为二进制文件名中的 OS/ARCH 部分
# 二进制命名约定: {cmd}-{os}-{arch}[.exe]
platform_to_pattern() {
  local platform="$1"
  local os="${platform%-*}"
  local arch="${platform#*-}"
  echo "${os}-${arch}"
}

# ============================================================
# 按平台打包
# ============================================================
for platform in "${PLATFORMS[@]}"; do
  ARCHIVE_NAME="tsign-skills-${VERSION}-${platform}.tar.gz"
  echo "  → 打包 ${platform}..."

  # 创建临时 staging 目录
  STAGING="$(mktemp -d)"
  trap "rm -rf '$STAGING'" EXIT

  PATTERN="$(platform_to_pattern "$platform")"
  SUFFIX="$(bin_suffix "$platform")"

  # 复制安装脚本到 staging 根目录
  cp "${TOOLKIT_DIR}/install.sh" "${STAGING}/"
  chmod +x "${STAGING}/install.sh"
  # Windows 平台额外打包 PowerShell 安装脚本
  if [[ "$platform" == windows-* ]] && [ -f "${TOOLKIT_DIR}/install.ps1" ]; then
    cp "${TOOLKIT_DIR}/install.ps1" "${STAGING}/"
  fi

  # 复制每个 skill
  for skill in "${SKILLS[@]}"; do
    SKILL_SRC="${SKILLS_DIR}/${skill}"
    SKILL_DST="${STAGING}/${skill}"
    mkdir -p "${SKILL_DST}/scripts/bin"

    # 复制 SKILL.md
    cp "${SKILL_SRC}/SKILL.md" "${SKILL_DST}/"

    # 复制 run.sh
    cp "${SKILL_SRC}/scripts/run.sh" "${SKILL_DST}/scripts/"
    chmod +x "${SKILL_DST}/scripts/run.sh"
    # Windows 平台额外打包 PowerShell 运行脚本
    if [[ "$platform" == windows-* ]] && [ -f "${SKILL_SRC}/scripts/run.ps1" ]; then
      cp "${SKILL_SRC}/scripts/run.ps1" "${SKILL_DST}/scripts/"
    fi

    # 只复制当前平台的二进制文件
    # 二进制文件按 bin/{os}-{arch}/ 子目录组织
    BIN_PLATFORM_DIR="${SKILL_SRC}/scripts/bin/${PATTERN}"
    if [ -d "$BIN_PLATFORM_DIR" ]; then
      for bin_file in "${BIN_PLATFORM_DIR}"/*; do
        [ -f "$bin_file" ] || continue
        cp "$bin_file" "${SKILL_DST}/scripts/bin/"
        chmod +x "${SKILL_DST}/scripts/bin/$(basename "$bin_file")"
      done

      # 统计复制的二进制数
      bin_count=$(find "${SKILL_DST}/scripts/bin" -type f 2>/dev/null | wc -l | tr -d ' ')
      if [ "$bin_count" -eq 0 ]; then
        echo "    ⚠ ${skill}: 未找到 ${platform} 的二进制文件"
      fi
    else
      echo "    ⚠ ${skill}: 未找到 ${platform} 的二进制目录 (${BIN_PLATFORM_DIR})"
    fi
  done

  # 打包
  tar -czf "${OUTPUT_DIR}/${ARCHIVE_NAME}" -C "$STAGING" .

  # 计算大小
  size=$(du -h "${OUTPUT_DIR}/${ARCHIVE_NAME}" | cut -f1 | tr -d ' ')
  echo "    ✅ ${ARCHIVE_NAME} (${size})"

  # 清理 staging
  rm -rf "$STAGING"
  # 重置 trap 以避免重复清理
  trap - EXIT
done

echo ""
echo "📋 生成校验文件..."

# 将 OUTPUT_DIR 转为绝对路径（cd 后仍可正确引用）
OUTPUT_DIR="$(cd "$OUTPUT_DIR" && pwd)"

# 生成 SHA256 校验文件
cd "$OUTPUT_DIR"
if command -v sha256sum &>/dev/null; then
  sha256sum tsign-skills-*.tar.gz > checksums-sha256.txt
elif command -v shasum &>/dev/null; then
  shasum -a 256 tsign-skills-*.tar.gz > checksums-sha256.txt
fi

echo ""
echo "✅ 打包完成！"
echo ""
echo "产物清单:"
ls -lh "${OUTPUT_DIR}"/tsign-skills-*.tar.gz 2>/dev/null | awk '{print "  " $NF " (" $5 ")"}'
echo ""
echo "校验文件: ${OUTPUT_DIR}/checksums-sha256.txt"
echo ""
echo "tarball 内部结构:"
echo "  tsign-skills-${VERSION}-{platform}.tar.gz"
echo "  ├── install.sh"
echo "  ├── install.ps1          (仅 Windows 平台)"
echo "  ├── contract-atoms/"
echo "  │   ├── SKILL.md"
echo "  │   └── scripts/"
echo "  │       ├── run.sh"
echo "  │       ├── run.ps1      (仅 Windows 平台)"
echo "  │       └── bin/   (当前平台二进制)"
echo "  ├── contract-review/"
echo "  │   └── ..."
echo "  └── contract-comparison/"
echo "      └── ..."
