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

TOOLKIT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$TOOLKIT_DIR")"
SKILLS_SRC="${PROJECT_ROOT}/skills"
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

echo "📦 安装腾讯电子签合同智能 Skills"
echo "   目标工具: ${TOOL_NAME}"
echo "   安装目录: ${DEST_DIR}"
echo ""

for skill in "${SKILLS[@]}"; do
  SRC="${SKILLS_SRC}/${skill}"
  if [ ! -d "$SRC" ]; then
    echo "⚠ 跳过 ${skill}（目录不存在）"
    continue
  fi

  SKILL_DEST="${DEST_DIR}/${skill}"
  echo "  → 安装 ${skill}..."
  mkdir -p "$SKILL_DEST"
  cp -r "$SRC/SKILL.md" "$SKILL_DEST/"
  cp -r "$SRC/scripts" "$SKILL_DEST/"

  # 确保二进制文件可执行
  find "$SKILL_DEST/scripts/bin" -type f 2>/dev/null | xargs chmod +x 2>/dev/null || true
  chmod +x "$SKILL_DEST/scripts/run.sh"
done

echo ""
echo "✅ 安装完成！已安装 ${#SKILLS[@]} 个 skills 到 ${DEST_DIR}"
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
