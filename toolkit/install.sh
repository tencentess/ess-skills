#!/usr/bin/env bash
# 腾讯电子签合同智能 Skills — 一键安装脚本
# 将 skills 安装到 AI 编码工具的 skills 目录
#
# 用法:
#   bash install.sh --tool=<claude|codebuddy|opencode> [--target=project|personal]
#
# 选项:
#   --tool=claude      安装到 Claude Code skills 目录 (默认)
#   --tool=codebuddy   安装到 CodeBuddy skills 目录
#   --tool=opencode    安装到 OpenCode skills 目录
#   --target=project   安装到当前项目目录 (默认)
#   --target=personal  安装到个人/全局目录
set -euo pipefail

# 脚本所在目录（本地开发时在 toolkit/，打包后在 tarball 根目录）
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 自动检测 skills 源目录：
#   打包产物：install.sh 与 contract-*/  同级
#   本地开发：install.sh 在 toolkit/，skills 在 ../skills/
if [ -d "${SCRIPT_DIR}/contract-atoms" ]; then
  SKILLS_SRC="$SCRIPT_DIR"
elif [ -d "${SCRIPT_DIR}/../skills" ]; then
  SKILLS_SRC="$(cd "${SCRIPT_DIR}/../skills" && pwd)"
else
  echo "❌ 找不到 skills 目录，请在安装包解压目录或项目根目录下运行" >&2
  exit 1
fi
TOOL="codebuddy"
TARGET="project"

for arg in "$@"; do
  case "$arg" in
    --tool=*) TOOL="${arg#*=}" ;;
    --target=*) TARGET="${arg#*=}" ;;
    --help|-h)
      echo "用法: bash install.sh --tool=<claude|codebuddy|opencode> [--target=project|personal]"
      echo ""
      echo "选项:"
      echo "  --tool=codebuddy   安装到 CodeBuddy skills 目录 (默认)"
      echo "  --tool=claude      安装到 Claude Code skills 目录"
      echo "  --tool=opencode    安装到 OpenCode skills 目录"
      echo "  --target=project   安装到当前项目目录 (默认)"
      echo "  --target=personal  安装到个人/全局目录"
      echo ""
      echo "示例:"
      echo "  bash install.sh --tool=claude --target=project"
      echo "  bash install.sh --tool=codebuddy --target=personal"
      echo "  bash install.sh --tool=opencode --target=project"
      exit 0
      ;;
    *)
      echo "❌ 未知参数: $arg" >&2
      echo "   运行 bash install.sh --help 查看用法" >&2
      exit 1
      ;;
  esac
done

# 根据工具和目标确定安装路径
case "$TOOL" in
  claude)
    TOOL_DIR=".claude/skills"
    PERSONAL_DIR="${HOME}/.claude/skills"
    TOOL_NAME="Claude Code"
    ;;
  codebuddy)
    TOOL_DIR=".codebuddy/skills"
    PERSONAL_DIR="${HOME}/.codebuddy/skills"
    TOOL_NAME="CodeBuddy"
    ;;
  opencode)
    TOOL_DIR=".opencode/skills"
    PERSONAL_DIR="${HOME}/.opencode/skills"
    TOOL_NAME="OpenCode"
    ;;
  *)
    echo "❌ 不支持的工具: $TOOL (可选: claude, codebuddy, opencode)" >&2
    exit 1
    ;;
esac

case "$TARGET" in
  project)
    DEST_DIR="$(pwd)/${TOOL_DIR}"
    ;;
  personal)
    DEST_DIR="$PERSONAL_DIR"
    ;;
  *)
    echo "❌ 无效的 target: $TARGET (可选: project, personal)" >&2
    exit 1
    ;;
esac

SKILLS=(contract-atoms contract-review contract-comparison)

# 检测当前平台
detect_platform() {
  local os arch
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"
  case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
  esac
  echo "${os}-${arch}"
}

PLATFORM="$(detect_platform)"

echo "📦 安装腾讯电子签合同智能 Skills"
echo "   目标工具: ${TOOL_NAME}"
echo "   安装目录: ${DEST_DIR}"
echo "   当前平台: ${PLATFORM}"
echo ""

INSTALLED=0
for skill in "${SKILLS[@]}"; do
  SRC="${SKILLS_SRC}/${skill}"
  if [ ! -d "$SRC" ]; then
    echo "⚠ 跳过 ${skill}（目录不存在）"
    continue
  fi

  SKILL_DEST="${DEST_DIR}/${skill}"
  echo "  → 安装 ${skill}..."
  mkdir -p "${SKILL_DEST}/scripts/bin"

  # 复制 SKILL.md
  cp "$SRC/SKILL.md" "$SKILL_DEST/"

  # 复制 run 脚本
  cp "$SRC/scripts/run.sh" "$SKILL_DEST/scripts/"
  chmod +x "$SKILL_DEST/scripts/run.sh"
  [ -f "$SRC/scripts/run.ps1" ] && cp "$SRC/scripts/run.ps1" "$SKILL_DEST/scripts/"

  # 复制二进制文件（区分 tarball 平面结构和本地开发的平台子目录结构）
  BIN_COPIED=0
  if [ -d "$SRC/scripts/bin/${PLATFORM}" ]; then
    # 本地开发：bin 按平台子目录组织 (bin/darwin-arm64/xxx)
    for f in "$SRC/scripts/bin/${PLATFORM}"/*; do
      [ -f "$f" ] || continue
      cp "$f" "$SKILL_DEST/scripts/bin/"
      chmod +x "$SKILL_DEST/scripts/bin/$(basename "$f")"
      BIN_COPIED=$((BIN_COPIED + 1))
    done
  elif ls "$SRC/scripts/bin/"* &>/dev/null; then
    # tarball：bin 是平面结构 (bin/xxx)
    for f in "$SRC/scripts/bin/"*; do
      [ -f "$f" ] || continue
      cp "$f" "$SKILL_DEST/scripts/bin/"
      chmod +x "$SKILL_DEST/scripts/bin/$(basename "$f")"
      BIN_COPIED=$((BIN_COPIED + 1))
    done
  fi

  if [ "$BIN_COPIED" -eq 0 ]; then
    echo "    ⚠ 未找到 ${PLATFORM} 的二进制文件（请先运行 make release）"
  else
    echo "    ✅ 复制 ${BIN_COPIED} 个二进制文件"
    INSTALLED=$((INSTALLED + 1))
  fi
done

echo ""
if [ "$INSTALLED" -eq 0 ]; then
  echo "⚠ 安装完成但缺少二进制文件。请先编译："
  echo ""
  echo "  cd $(cd "$SCRIPT_DIR" && pwd) && make release"
  echo ""
  echo "然后重新运行 install.sh"
else
  echo "✅ 安装完成！已安装 ${INSTALLED}/${#SKILLS[@]} 个 skills 到 ${DEST_DIR}"
fi
echo ""
echo "下一步：配置凭证（任选一种方式）："
echo ""
echo "  方式一（推荐）：创建 ~/.tsign/config.yaml"
echo "    credentials:"
echo "      secret_id: \"AKIDxxxxxxxx\""
echo "      secret_key: \"xxxxxxxx\""
echo "    operator:"
echo "      user_id: \"yDwJxxx\""
echo "    env: \"online\""
echo ""
echo "  方式二：设置环境变量（适合 CI/CD）"
echo "    export TENCENTCLOUD_SECRET_ID=\"AKIDxxxxxxxx\""
echo "    export TENCENTCLOUD_SECRET_KEY=\"xxxxxxxx\""
echo "    export ESS_OPERATOR_ID=\"yDwJxxx\""
