# 合同智能 Skills

[English](./README_EN.md) | 中文

腾讯电子签合同智能 Skill 集合 — 为 Agent（CodeBuddy / Claude Code ）提供合同审查、合同对比能力。

## Skills 列表

| Skill | 说明 | 详情 |
|-------|------|------|
| [contract-review](./skills/contract-review/SKILL.md) | 合同审查工作流：上传 PDF → AI 风险识别 → 输出风险报告 | 一键完成端到端审查 |
| [contract-comparison](./skills/contract-comparison/SKILL.md) | 合同对比工作流：上传两份文件 → 差异分析 → 输出对比报告 | 支持 PDF、Word 等格式 |
| [contract-atoms](./skills/contract-atoms/SKILL.md) | 7 个原子操作命令 | 适合精细控制流程的场景 |

## 安装

### CodeBuddy

```bash
bash toolkit/install.sh --tool=codebuddy
```

### Claude Code

```bash
bash toolkit/install.sh --tool=claude
```



> 默认安装到当前项目目录，加 `--target=personal` 安装到个人全局目录。

### 远程一键安装

无需克隆仓库，直接下载安装：

**macOS / Linux**

```bash
curl -fsSL https://raw.githubusercontent.com/tencentess/ess-skills/main/toolkit/install-remote.sh | bash
```

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/tencentess/ess-skills/main/toolkit/install-remote.ps1 -OutFile install-remote.ps1; powershell -ExecutionPolicy Bypass -File .\install-remote.ps1
```

指定参数：

```powershell
powershell -ExecutionPolicy Bypass -File .\install-remote.ps1 -Tool claude -Target personal
```

## 凭证配置

使用前需配置腾讯电子签凭证，共需要三个参数：

| 参数 | 说明 | 获取方式 |
|------|------|----------|
| `secret_id` | 腾讯云 API 密钥 SecretId | 前往 [腾讯云 API 密钥管理](https://console.cloud.tencent.com/cam/capi) 创建或查看 |
| `secret_key` | 腾讯云 API 密钥 SecretKey | 同上，与 SecretId 成对获取 |
| `user_id` | 经办人/操作人编号 (UserId) | 登录 [电子签控制台](https://qian.tencent.com) → 组织管理 → 组织架构，查询员工 UserId |

> **密钥安全提示**：SecretId / SecretKey 是企业身份凭证，请妥善保管，切勿泄露或提交到代码仓库。
>
> 详细说明：[密钥获取](https://qian.tencent.com/developers/company/online_env_integration#2%E8%8E%B7%E5%8F%96%E5%AF%86%E9%92%A5secretid%E5%92%8Csecretkey%E7%BA%BF%E4%B8%8A%E7%8E%AF%E5%A2%83) | [UserId 获取](https://qian.tencent.com/developers/company/common_params#%E4%B8%80-%E7%BB%8F%E5%8A%9E%E4%BA%BA%E6%93%8D%E4%BD%9C%E4%BA%BA%E7%BC%96%E5%8F%B7-userid-%E8%8E%B7%E5%8F%96)

### 方式一：配置文件（推荐）

**macOS / Linux**：创建 `~/.tsign/config.yaml`

```bash
mkdir -p ~/.tsign && cat > ~/.tsign/config.yaml << 'EOF'
credentials:
  secret_id: "AKIDxxxxxxxx"
  secret_key: "xxxxxxxx"
operator:
  user_id: "yDwJxxx"
env: "online"
EOF
```

**Windows (PowerShell)**：创建 `%USERPROFILE%\.tsign\config.yaml`

```powershell
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.tsign" | Out-Null
@"
credentials:
  secret_id: "AKIDxxxxxxxx"
  secret_key: "xxxxxxxx"
operator:
  user_id: "yDwJxxx"
env: "online"
"@ | Set-Content "$env:USERPROFILE\.tsign\config.yaml" -Encoding UTF8
```

### 方式二：环境变量

**macOS / Linux (Bash)**

```bash
export TENCENTCLOUD_SECRET_ID="AKIDxxxxxxxx"
export TENCENTCLOUD_SECRET_KEY="xxxxxxxx"
export ESS_OPERATOR_ID="yDwJxxx"
```

**Windows (PowerShell)**

```powershell
$env:TENCENTCLOUD_SECRET_ID="AKIDxxxxxxxx"
$env:TENCENTCLOUD_SECRET_KEY="xxxxxxxx"
$env:ESS_OPERATOR_ID="yDwJxxx"
```

> 运行命令时若未检测到凭证，会自动引导创建。详见 [凭证配置说明](./toolkit/references/credentials.md)。

安装完成后，在 AI 助手中直接对话即可使用，例如"帮我审查这份合同"、"对比这两份合同的差异"。具体命令和参数详见各 Skill 的 SKILL.md。

## 目录结构

```
ess-skills/
├── skills/                  # Skills（SKILL.md + scripts/ + 预编译二进制）
│   ├── contract-atoms/
│   ├── contract-review/
│   └── contract-comparison/
├── toolkit/                 # 源码 & 开发工具（详见 toolkit/README.md）
└── Makefile
```

## 开发

构建、测试、架构设计等开发相关信息请参见 [toolkit/README.md](./toolkit/README.md)。

## License

MIT
