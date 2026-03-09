# Contract Intelligence Skills

English | [中文](./README.md)

Tencent E-Sign Contract Intelligence Skill collection — provides contract review, comparison, and atomic operations for Agents(CodeBuddy / Claude Code).

## Skills

| Skill | Description | Details |
|-------|-------------|---------|
| [contract-review](./skills/contract-review/SKILL.md) | Contract review workflow: Upload PDF → AI risk analysis → Risk report | End-to-end one-click review |
| [contract-comparison](./skills/contract-comparison/SKILL.md) | Contract comparison workflow: Upload two files → Diff analysis → Comparison report | Supports PDF, Word, etc. |
| [contract-atoms](./skills/contract-atoms/SKILL.md) | 7 atomic CLI commands | For fine-grained control |

## Installation

### CodeBuddy

```bash
bash toolkit/install.sh --tool=codebuddy
```

### Claude Code

```bash
bash toolkit/install.sh --tool=claude
```


> Installs to the current project by default. Add `--target=personal` for global installation.

### Remote One-Click Install

No need to clone the repo — download and install directly:

**macOS / Linux**

```bash
curl -fsSL https://raw.githubusercontent.com/tencentess/ess-skills/main/toolkit/install-remote.sh | bash
```

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/tencentess/ess-skills/main/toolkit/install-remote.ps1 -OutFile install-remote.ps1; powershell -ExecutionPolicy Bypass -File .\install-remote.ps1
```

With parameters:

```powershell
powershell -ExecutionPolicy Bypass -File .\install-remote.ps1 -Tool claude -Target personal
```

## Credential Setup

Three parameters are required to configure Tencent E-Sign credentials:

| Parameter | Description | How to Obtain |
|-----------|-------------|---------------|
| `secret_id` | Tencent Cloud API SecretId | Go to [Tencent Cloud API Key Management](https://console.cloud.tencent.com/cam/capi) to create or view |
| `secret_key` | Tencent Cloud API SecretKey | Same as above, obtained as a pair with SecretId |
| `user_id` | Operator UserId | Log in to [E-Sign Console](https://qian.tencent.com) → Organization → Members, find the UserId |

> **Security Note**: SecretId / SecretKey are enterprise identity credentials. Keep them safe and never commit them to a repository.
>
> Reference: [Key Acquisition](https://qian.tencent.com/developers/company/online_env_integration#2%E8%8E%B7%E5%8F%96%E5%AF%86%E9%92%A5secretid%E5%92%8Csecretkey%E7%BA%BF%E4%B8%8A%E7%8E%AF%E5%A2%83) | [UserId Acquisition](https://qian.tencent.com/developers/company/common_params#%E4%B8%80-%E7%BB%8F%E5%8A%9E%E4%BA%BA%E6%93%8D%E4%BD%9C%E4%BA%BA%E7%BC%96%E5%8F%B7-userid-%E8%8E%B7%E5%8F%96)

### Option 1: Configuration File (Recommended)

**macOS / Linux**: Create `~/.tsign/config.yaml`

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

**Windows (PowerShell)**: Create `%USERPROFILE%\.tsign\config.yaml`

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

### Option 2: Environment Variables

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

> If no credentials are detected, an interactive setup wizard will guide you. See [Credential Setup](./toolkit/references/credentials.md) for details.

Once installed, simply chat with your AI assistant — e.g. "review this contract" or "compare these two contracts". See each Skill's SKILL.md for command details and parameters.

## Directory Structure

```
ess-skills/
├── skills/                  # Skills (SKILL.md + scripts/ + pre-built binaries)
│   ├── contract-atoms/
│   ├── contract-review/
│   └── contract-comparison/
├── toolkit/                 # Source code & dev tools (see toolkit/README.md)
└── Makefile
```

## Development

For building, testing, architecture details, see [toolkit/README.md](./toolkit/README.md).

## License

MIT
