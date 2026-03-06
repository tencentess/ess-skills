# Toolkit — 开发指南

合同智能 Skills 的源码、构建工具和开发文档。

## 三层架构

```
┌─────────────────────────────────────────────┐
│  Layer 2: 端到端工作流                        │
│  contract-review / contract-comparison       │
│  一键完成: 上传→创建任务→轮询等待→获取结果→导出   │
└─────────────────┬───────────────────────────┘
                  │ 组合编排
┌─────────────────▼───────────────────────────┐
│  Layer 1: 原子操作                            │
│  contract-atoms (7 个独立 CLI 命令)            │
└─────────────────┬───────────────────────────┘
                  │ 引用
┌─────────────────▼───────────────────────────┐
│  Layer 0: 共享基础设施                         │
│  foundation (Go 共享库)                       │
│  client / config / output / poller            │
└─────────────────────────────────────────────┘
```

## 目录结构

```
toolkit/
├── foundation/              # Layer 0: 共享基础设施（Go 库）
│   ├── client/              #   API 客户端封装
│   ├── config/              #   凭证 & 配置加载
│   ├── output/              #   统一 JSON 输出
│   └── poller/              #   异步任务轮询（指数退避）
├── contract-atoms/          # Layer 1: 7 个原子操作 CLI
├── contract-review/         # Layer 2: 合同审查端到端工作流
├── contract-comparison/     # Layer 2: 合同对比端到端工作流
├── scripts/                 # 构建 & 打包脚本
├── references/              # 参考文档（凭证配置说明等）
├── test/                    # 测试用例 & 样例文件
├── install.sh               # 安装脚本（macOS/Linux）
├── install.ps1              # 安装脚本（Windows）
├── install-remote.sh        # 远程一键安装（macOS/Linux）
├── install-remote.ps1       # 远程一键安装（Windows）
└── Makefile                 # 构建入口
```

## 构建

```bash
# 全量编译（6 平台交叉编译）
make build-all

# 构建并打包到 skills/
make release

# 按平台拆分打包为 tar.gz（输出到 dist/）
make package VERSION=1.0.0

# 清理构建产物
make clean

# 运行测试
make test
```

### 目标平台

| OS | Arch |
|----|------|
| darwin | amd64, arm64 |
| linux | amd64, arm64 |
| windows | amd64, arm64 |

## 构建流程

1. `make build-all` — 对每个子模块执行 `go mod tidy` + 6 平台交叉编译
2. `make release` — 在 `build-all` 基础上，将 SKILL.md + scripts/ + 预编译二进制同步到 `skills/` 目录
3. `make package` — 在 `release` 基础上，按平台拆分打包为 `dist/tsign-skills-{version}-{os}-{arch}.tar.gz`

## 输出格式

所有 CLI 命令统一输出 JSON 到 stdout：

```json
// 成功
{"success": true, "data": {...}, "error": null}

// 失败
{"success": false, "data": null, "error": {"code": "xxx", "message": "xxx"}}
```

## 凭证配置

详见 [references/credentials.md](./references/credentials.md)。
