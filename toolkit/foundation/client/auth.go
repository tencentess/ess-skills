package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tencentess/ess-skills/toolkit/foundation/config"
	"golang.org/x/term"
)

// CLIFlags 命令行凭证参数
type CLIFlags struct {
	SecretID   string
	SecretKey  string
	OperatorID string
	Profile    string
	Env        string // test / online，默认 online
}

// ResolvedCredentials 已解析的完整凭证
type ResolvedCredentials struct {
	SecretID   string
	SecretKey  string
	OperatorID string
	Env        string // test / online
}

// LoadCredentials 按优先级链加载凭证：CLI参数 > 环境变量 > 配置文件 > 交互式
func LoadCredentials(cli *CLIFlags) (*ResolvedCredentials, error) {
	// 优先级 1: 命令行参数
	if cli != nil && cli.SecretID != "" && cli.SecretKey != "" && cli.OperatorID != "" {
		env := cli.Env
		if env == "" {
			env = envOrDefault("ESS_ENV", "online")
		}
		return &ResolvedCredentials{
			SecretID:   cli.SecretID,
			SecretKey:  cli.SecretKey,
			OperatorID: cli.OperatorID,
			Env:        env,
		}, nil
	}

	// 优先级 2: 环境变量（兼容腾讯云 SDK 标准命名）
	envID := firstNonEmpty(os.Getenv("TENCENTCLOUD_SECRET_ID"), os.Getenv("TENCENT_SECRET_ID"))
	envKey := firstNonEmpty(os.Getenv("TENCENTCLOUD_SECRET_KEY"), os.Getenv("TENCENT_SECRET_KEY"))
	envOp := os.Getenv("ESS_OPERATOR_ID")
	if envID != "" && envKey != "" && envOp != "" {
		env := envOrDefault("ESS_ENV", "online")
		return &ResolvedCredentials{
			SecretID: envID, SecretKey: envKey, OperatorID: envOp, Env: env,
		}, nil
	}

	// 优先级 3: 配置文件
	profile := ""
	if cli != nil {
		profile = cli.Profile
	}
	cfg, err := config.Load(profile)
	if err == nil && cfg.Credentials.SecretID != "" && cfg.Credentials.SecretKey != "" && cfg.Operator.UserID != "" {
		env := cfg.Env
		if env == "" {
			env = "online"
		}
		return &ResolvedCredentials{
			SecretID:   cfg.Credentials.SecretID,
			SecretKey:  cfg.Credentials.SecretKey,
			OperatorID: cfg.Operator.UserID,
			Env:        env,
		}, nil
	}

	// 优先级 4: 交互式引导（仅 TTY）
	if isTerminal() {
		return interactiveSetup()
	}

	return nil, fmt.Errorf("未找到凭证。请创建配置文件 (%s) 或设置环境变量。\n\n"+
		"📖 获取说明:\n"+
		"  • SecretId/SecretKey 获取: https://qian.tencent.com/developers/company/online_env_integration#2\n"+
		"  • 经办人(OperatorId) 获取: https://qian.tencent.com/developers/company/common_params#%s\n\n"+
		"💡 快速配置（推荐创建 %s）:\n"+
		"  credentials:\n"+
		"    secret_id: \"AKIDxxxxxxxx\"\n"+
		"    secret_key: \"xxxxxxxx\"\n"+
		"  operator:\n"+
		"    user_id: \"yDwJxxx\"\n"+
		"  env: \"online\"\n\n"+
		"  或设置环境变量（适合 CI/CD）:\n"+
		"  export TENCENTCLOUD_SECRET_ID=\"AKIDxxxxxxxx\"\n"+
		"  export TENCENTCLOUD_SECRET_KEY=\"xxxxxxxx\"\n"+
		"  export ESS_OPERATOR_ID=\"yDwJxxx\"",
		config.GetConfigPath(),
		"%E4%B8%80-%E7%BB%8F%E5%8A%9E%E4%BA%BA%E6%93%8D%E4%BD%9C%E4%BA%BA%E7%BC%96%E5%8F%B7-userid-%E8%8E%B7%E5%8F%96",
		config.GetConfigPath())
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func isTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func interactiveSetup() (*ResolvedCredentials, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintln(os.Stderr, "⚠ 未检测到凭证配置。开始交互式配置...")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "📖 获取说明:")
	fmt.Fprintln(os.Stderr, "  • SecretId/SecretKey: https://qian.tencent.com/developers/company/online_env_integration#2")
	fmt.Fprintln(os.Stderr, "  • 经办人(OperatorId): https://qian.tencent.com/developers/company/common_params")
	fmt.Fprintln(os.Stderr)

	fmt.Fprint(os.Stderr,    "Secret ID: ")
	sid, _ := reader.ReadString('\n')
	sid = strings.TrimSpace(sid)

	fmt.Fprint(os.Stderr, "Secret Key: ")
	skeyBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("读取 Secret Key 失败: %w", err)
	}
	skey := string(skeyBytes)
	fmt.Fprintln(os.Stderr)

	fmt.Fprint(os.Stderr, "Operator UserId: ")
	opid, _ := reader.ReadString('\n')
	opid = strings.TrimSpace(opid)

	fmt.Fprint(os.Stderr, "环境 (test/online) [online]: ")
	env, _ := reader.ReadString('\n')
	env = strings.TrimSpace(env)
	if env == "" {
		env = "online"
	}

	cred := &ResolvedCredentials{SecretID: sid, SecretKey: skey, OperatorID: opid, Env: env}

	// 询问是否保存
	fmt.Fprint(os.Stderr, "\n是否保存到配置文件? [Y/n]: ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer == "" || answer == "y" || answer == "yes" {
		cfg := &config.Config{
			Credentials: config.Credentials{SecretID: sid, SecretKey: skey},
			Operator:    config.Operator{UserID: opid},
			Env:         env,
		}
		if saveErr := config.Save(cfg); saveErr != nil {
			fmt.Fprintf(os.Stderr, "⚠ 保存配置失败: %v\n", saveErr)
		} else {
			fmt.Fprintf(os.Stderr, "✅ 配置已保存到 %s\n", config.GetConfigPath())
		}
	}
	return cred, nil
}
