## 前置配置

使用前需配置腾讯电子签凭证，**任选一种方式**：

### 方式一：配置文件（推荐）

创建 `~/.tsign/config.yaml`（运行命令时若未检测到凭证会自动引导创建）：

```yaml
# 默认凭证
credentials:
  secret_id: "AKIDxxxxxxxx"
  secret_key: "xxxxxxxx"

# 经办人
operator:
  user_id: "yDwJxxx"

# 环境: test / online (默认 online)
env: "online"

# 可选: 多 profile 配置（通过 --profile=dev 切换）
profiles:
  dev:
    credentials:
      secret_id: "AKIDdev_xxx"
      secret_key: "dev_xxx"
    operator:
      user_id: "yDwJdev"
    env: "test"
  prod:
    credentials:
      secret_id: "AKIDprod_xxx"
      secret_key: "prod_xxx"
    operator:
      user_id: "yDwJprod"
    env: "online"
```

> **安全提示**: 配置文件权限自动设置为 `0600`（仅本人可读写）。也可通过环境变量 `TSIGN_CONFIG_PATH` 自定义路径。

### 方式二：环境变量（适合 CI/CD 等自动化场景）

**macOS / Linux (bash/zsh):**
```bash
export TENCENTCLOUD_SECRET_ID="AKIDxxxxxxxx"
export TENCENTCLOUD_SECRET_KEY="xxxxxxxx"
export ESS_OPERATOR_ID="yDwJxxx"
```

**Windows (PowerShell):**
```powershell
$env:TENCENTCLOUD_SECRET_ID = "AKIDxxxxxxxx"
$env:TENCENTCLOUD_SECRET_KEY = "xxxxxxxx"
$env:ESS_OPERATOR_ID = "yDwJxxx"
```

> **注意**: 环境变量优先级高于配置文件。若两者同时存在，环境变量会覆盖配置文件中的值。

### 凭证覆盖参数

所有命令都支持以下参数覆盖凭证：
- `--secret-id`: 覆盖 SecretId
- `--secret-key`: 覆盖 SecretKey
- `--operator-id`: 覆盖经办人 UserId
- `--profile`: 指定配置文件 profile 名称
- `--env`: 指定环境 test/online (默认 online)
