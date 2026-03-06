#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"  # toolkit/

echo "🔨 Building all contract intelligence skills..."
echo ""

cd "$ROOT_DIR"

# 确保 foundation 依赖就绪
echo "=== Preparing foundation ==="
cd foundation && go mod tidy && cd ..

# 编译所有子项目
for dir in contract-atoms contract-review contract-comparison; do
    echo ""
    echo "=== Building $dir ==="
    cd "$dir"
    go mod tidy
    make build-all
    cd ..
done

echo ""
echo "🎉 All builds complete!"
echo ""

# 统计产出
for dir in contract-atoms contract-review contract-comparison; do
    count=$(find "$dir/bin" -type f 2>/dev/null | wc -l | tr -d ' ')
    echo "  $dir: $count binaries"
done
