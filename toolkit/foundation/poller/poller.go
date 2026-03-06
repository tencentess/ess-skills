package poller

import (
	"context"
	"fmt"
	"os"
	"time"
)

// PollConfig 轮询配置
type PollConfig struct {
	InitialInterval time.Duration // 首次轮询间隔
	MaxInterval     time.Duration // 最大轮询间隔
	Timeout         time.Duration // 总超时
	BackoffFactor   float64       // 退避倍数
	ResetInterval   time.Duration // 到达 MaxInterval 后重置为此值（实现周期性退避）
}

// DefaultPollConfig 默认轮询配置
var DefaultPollConfig = PollConfig{
	InitialInterval: 1 * time.Second,
	MaxInterval:     15 * time.Second,
	Timeout:         10 * time.Minute,
	BackoffFactor:   1.5,
	ResetInterval:   1 * time.Second, // 到达 MaxInterval 后重置为此值
}

// PollResult 单次轮询结果
type PollResult struct {
	Done bool
	Data interface{}
	Err  error
}

// Poll 通用轮询函数，fn 返回 PollResult 表示是否完成
func Poll(ctx context.Context, cfg PollConfig, fn func() PollResult) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	interval := cfg.InitialInterval
	pollCount := 0
	startTime := time.Now()
	for {
		result := fn()
		if result.Err != nil {
			if pollCount > 0 {
				fmt.Fprintln(os.Stderr) // 换行，结束进度行
			}
			return nil, result.Err
		}
		if result.Done {
			if pollCount > 0 {
				fmt.Fprintln(os.Stderr) // 换行，结束进度行
			}
			return result.Data, nil
		}

		pollCount++
		elapsed := time.Since(startTime).Truncate(time.Second)
		secs := int(interval.Seconds())
		if secs < 1 {
			secs = 1
		}
		// 使用 \r 回到行首覆盖输出，避免刷屏
		fmt.Fprintf(os.Stderr, "\r⏳ 任务执行中... 已等待 %v，%d秒后第%d次重试", elapsed, secs, pollCount)

		select {
		case <-ctx.Done():
			fmt.Fprintln(os.Stderr) // 换行
			return nil, fmt.Errorf("轮询超时（已等待 %v）", cfg.Timeout)
		case <-time.After(interval):
		}

		// 指数退避，到达 MaxInterval 后重置为 ResetInterval 重新退避
		interval = time.Duration(float64(interval) * cfg.BackoffFactor)
		if interval >= cfg.MaxInterval {
			if cfg.ResetInterval > 0 {
				interval = cfg.ResetInterval
			} else {
				interval = cfg.MaxInterval
			}
		}
	}
}
