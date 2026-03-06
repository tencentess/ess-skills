package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tencentess/ess-skills/toolkit/foundation/client"
	"github.com/tencentess/ess-skills/toolkit/foundation/output"
	"github.com/tencentess/ess-skills/toolkit/foundation/poller"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ess "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ess/v20201111"
)

func main() {
	fs := flag.NewFlagSet("compare-workflow", flag.ExitOnError)

	var (
		secretID   string
		secretKey  string
		operatorID string
		profile    string
		env        string
		originFile string
		diffFile   string
		comment    string
		exportType int
		timeout    string
	)

	fs.StringVar(&secretID, "secret-id", "", "腾讯云 SecretId")
	fs.StringVar(&secretKey, "secret-key", "", "腾讯云 SecretKey")
	fs.StringVar(&operatorID, "operator-id", "", "经办人 UserId")
	fs.StringVar(&profile, "profile", "", "配置文件 profile 名称")
	fs.StringVar(&env, "env", "", "环境: test/online (默认 online)")
	fs.StringVar(&originFile, "origin-file", "", "原版 PDF 文件路径 (必填)")
	fs.StringVar(&diffFile, "diff-file", "", "新版 PDF 文件路径 (必填)")
	fs.StringVar(&comment, "comment", "", "备注")
	fs.IntVar(&exportType, "export", -1, "导出类型: -1=不导出 0=PDF报告 1=Excel明细 (默认-1)")
	fs.StringVar(&timeout, "timeout", "10m", "超时时间 (默认10m)")
	fs.Parse(os.Args[1:])

	if originFile == "" || diffFile == "" {
		output.PrintError("InvalidParameter", "--origin-file 和 --diff-file 为必填参数")
	}

	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		output.PrintErrorf("InvalidParameter", "无效的超时时间格式: %v", err)
	}

	// 1. 加载凭证
	fmt.Fprintln(os.Stderr, "🔐 加载凭证...")
	cred, err := client.LoadCredentials(&client.CLIFlags{
		SecretID: secretID, SecretKey: secretKey,
		OperatorID: operatorID, Profile: profile, Env: env,
	})
	if err != nil {
		output.PrintError("AuthFailure", err.Error())
	}

	essClient, err := client.NewEssClient(cred)
	if err != nil {
		output.PrintError("ClientError", err.Error())
	}

	// 2. 上传两份文件
	fmt.Fprintf(os.Stderr, "📤 上传原版文件: %s\n", originFile)
	originID, err := essClient.UploadLocalFile(originFile)
	if err != nil {
		output.PrintErrorf("UploadError", "原版文件上传失败: %v", err)
	}
	fmt.Fprintf(os.Stderr, "✅ 原版上传成功, ResourceId: %s\n", originID)

	fmt.Fprintf(os.Stderr, "📤 上传新版文件: %s\n", diffFile)
	diffID, err := essClient.UploadLocalFile(diffFile)
	if err != nil {
		output.PrintErrorf("UploadError", "新版文件上传失败: %v", err)
	}
	fmt.Fprintf(os.Stderr, "✅ 新版上传成功, ResourceId: %s\n", diffID)

	// 3. 创建对比任务
	fmt.Fprintln(os.Stderr, "📝 创建对比任务...")
	createReq := ess.NewCreateContractComparisonTaskRequest()
	createReq.Operator = essClient.Operator()
	createReq.OriginFileResourceId = common.StringPtr(originID)
	createReq.DiffFileResourceId = common.StringPtr(diffID)
	if comment != "" {
		createReq.Comment = common.StringPtr(comment)
	}

	createResp, err := essClient.SDKClient.CreateContractComparisonTask(createReq)
	if err != nil {
		output.PrintErrorf("CreateTaskError", "创建对比任务失败: %v", err)
	}
	if createResp.Response == nil || createResp.Response.TaskId == nil {
		output.PrintError("CreateTaskError", "创建任务成功但未返回 TaskId")
	}
	taskID := *createResp.Response.TaskId
	fmt.Fprintf(os.Stderr, "✅ 任务创建成功, TaskId: %s\n", taskID)

	// 4. 轮询等待
	fmt.Fprintln(os.Stderr, "⏳ 等待对比完成...")
	pollCfg := poller.DefaultPollConfig
	pollCfg.Timeout = timeoutDuration

	pollResult, err := poller.Poll(context.Background(), pollCfg, func() poller.PollResult {
		queryReq := ess.NewDescribeContractComparisonTaskRequest()
		queryReq.Operator = essClient.Operator()
		queryReq.TaskId = common.StringPtr(taskID)

		queryResp, err := essClient.SDKClient.DescribeContractComparisonTask(queryReq)
		if err != nil {
			return poller.PollResult{Err: fmt.Errorf("查询任务状态失败: %w", err)}
		}

		if queryResp.Response == nil || queryResp.Response.Status == nil {
			return poller.PollResult{Err: fmt.Errorf("查询任务返回空响应")}
		}

		status := *queryResp.Response.Status
		switch status {
		case 2: // 成功
			return poller.PollResult{Done: true, Data: queryResp.Response}
		case 3: // 失败
			msg := "对比任务执行失败"
			if queryResp.Response.Message != nil {
				msg = *queryResp.Response.Message
			}
			return poller.PollResult{Err: fmt.Errorf(msg)}
		default:
			return poller.PollResult{Done: false}
		}
	})
	if err != nil {
		output.PrintErrorf("PollError", "等待对比完成失败: %v", err)
	}

	fmt.Fprintln(os.Stderr, "✅ 对比完成!")

	compResp := pollResult.(*ess.DescribeContractComparisonTaskResponseParams)

	// 5. 获取对比详情（带 ShowDetail）
	fmt.Fprintln(os.Stderr, "📋 获取对比详情...")
	detailReq := ess.NewDescribeContractComparisonTaskRequest()
	detailReq.Operator = essClient.Operator()
	detailReq.TaskId = common.StringPtr(taskID)
	detailReq.ShowDetail = common.BoolPtr(true)

	detailResp, err := essClient.SDKClient.DescribeContractComparisonTask(detailReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠ 获取详情失败: %v，使用基础结果\n", err)
	} else if detailResp.Response != nil {
		compResp = detailResp.Response
	}

	// 6. 获取在线对比预览页面
	var webURL string
	fmt.Fprintln(os.Stderr, "🌐 获取在线对比预览页...")
	diffWebReq := ess.NewCreateContractDiffTaskWebUrlRequest()
	diffWebReq.Operator = essClient.Operator()
	diffWebReq.SkipFileUpload = common.BoolPtr(true)
	diffWebReq.OriginalFileResourceId = common.StringPtr(originID)
	diffWebReq.DiffFileResourceId = common.StringPtr(diffID)

	diffWebResp, err := essClient.SDKClient.CreateContractDiffTaskWebUrl(diffWebReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠ 获取预览页失败: %v\n", err)
	} else if diffWebResp.Response != nil {
		// 需要等待这个任务完成后再获取结果页面
		if diffWebResp.Response.TaskId != nil {
			diffTaskID := *diffWebResp.Response.TaskId
			fmt.Fprintln(os.Stderr, "⏳ 等待对比预览页生成...")

			// 轮询等待 DiffTaskWebUrl 完成
			webPollCfg := poller.DefaultPollConfig
			webPollCfg.Timeout = 3 * time.Minute

			webResult, webErr := poller.Poll(context.Background(), webPollCfg, func() poller.PollResult {
				descReq := ess.NewDescribeContractDiffTaskWebUrlRequest()
				descReq.Operator = essClient.Operator()
				descReq.TaskId = common.StringPtr(diffTaskID)

				descResp, err := essClient.SDKClient.DescribeContractDiffTaskWebUrl(descReq)
				if err != nil {
					return poller.PollResult{Err: fmt.Errorf("获取预览页失败: %w", err)}
				}
				if descResp.Response != nil && descResp.Response.WebUrl != nil && *descResp.Response.WebUrl != "" {
					return poller.PollResult{Done: true, Data: *descResp.Response.WebUrl}
				}
				return poller.PollResult{Done: false}
			})

			if webErr != nil {
				fmt.Fprintf(os.Stderr, "⚠ 获取预览页超时: %v\n", webErr)
			} else if url, ok := webResult.(string); ok {
				webURL = url
				fmt.Fprintln(os.Stderr, "✅ 预览页获取成功!")
			}
		}
		// CreateContractDiffTaskWebUrl 直接返回了 WebUrl
		if webURL == "" && diffWebResp.Response.WebUrl != nil && *diffWebResp.Response.WebUrl != "" {
			webURL = *diffWebResp.Response.WebUrl
			fmt.Fprintln(os.Stderr, "✅ 预览页获取成功!")
		}
	}

	// 7. 可选导出
	var exportURL string
	if exportType >= 0 {
		fmt.Fprintf(os.Stderr, "📥 导出结果 (类型: %d)...\n", exportType)
		exportReq := ess.NewExportContractComparisonTaskRequest()
		exportReq.Operator = essClient.Operator()
		exportReq.TaskId = common.StringPtr(taskID)
		exportReq.ExportType = common.Int64Ptr(int64(exportType))

		exportResp, err := essClient.SDKClient.ExportContractComparisonTask(exportReq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠ 导出失败: %v\n", err)
		} else if exportResp.Response != nil && exportResp.Response.ResourceUrl != nil {
			exportURL = *exportResp.Response.ResourceUrl
			fmt.Fprintln(os.Stderr, "✅ 导出成功!")
		}
	}

	// 8. 构建 Markdown 结果并输出
	md := renderMarkdown(compResp, taskID, filepath.Base(originFile), filepath.Base(diffFile), webURL, exportURL, exportType)

	// 直接将 Markdown 输出到 stderr，让终端/Agent 可以直接展示
	fmt.Fprintln(os.Stderr, "\n"+md)

	// 构建结构化输出（stdout JSON 供程序消费）
	finalResult := map[string]interface{}{
		"task_id":  taskID,
		"markdown": md,
	}
	if webURL != "" {
		finalResult["web_url"] = webURL
	}
	if exportURL != "" {
		finalResult["export_url"] = exportURL
	}
	// 添加统计数据
	stats := map[string]interface{}{}
	if compResp.TotalDiffCount != nil {
		stats["total"] = *compResp.TotalDiffCount
	}
	if compResp.AddDiffCount != nil {
		stats["add"] = *compResp.AddDiffCount
	}
	if compResp.ChangeDiffCount != nil {
		stats["change"] = *compResp.ChangeDiffCount
	}
	if compResp.DeleteDiffCount != nil {
		stats["delete"] = *compResp.DeleteDiffCount
	}
	finalResult["diff_stats"] = stats

	output.PrintSuccess(finalResult)
}

// renderMarkdown 将对比结果渲染为 Markdown 格式
func renderMarkdown(resp *ess.DescribeContractComparisonTaskResponseParams, taskID, originName, diffName, webURL, exportURL string, exportType int) string {
	var sb strings.Builder

	sb.WriteString("# 合同对比报告\n\n")

	// 基本信息
	sb.WriteString("## 基本信息\n\n")
	sb.WriteString("| 项目 | 内容 |\n")
	sb.WriteString("| --- | --- |\n")
	sb.WriteString(fmt.Sprintf("| 原版文件 | %s |\n", originName))
	sb.WriteString(fmt.Sprintf("| 新版文件 | %s |\n", diffName))
	sb.WriteString(fmt.Sprintf("| 任务ID | `%s` |\n", taskID))

	// 差异统计
	sb.WriteString("\n## 差异统计\n\n")
	sb.WriteString("| 类型 | 数量 |\n")
	sb.WriteString("| --- | --- |\n")
	if resp.TotalDiffCount != nil {
		sb.WriteString(fmt.Sprintf("| 差异总数 | **%d** |\n", *resp.TotalDiffCount))
	}
	if resp.AddDiffCount != nil {
		sb.WriteString(fmt.Sprintf("| 🟢 新增 | %d |\n", *resp.AddDiffCount))
	}
	if resp.ChangeDiffCount != nil {
		sb.WriteString(fmt.Sprintf("| 🟡 修改 | %d |\n", *resp.ChangeDiffCount))
	}
	if resp.DeleteDiffCount != nil {
		sb.WriteString(fmt.Sprintf("| 🔴 删除 | %d |\n", *resp.DeleteDiffCount))
	}
	sb.WriteString("\n")

	// 重要链接
	sb.WriteString("## 📌 重要链接\n\n")
	if webURL != "" {
		sb.WriteString(fmt.Sprintf("> **🌐 在线对比预览页面（可直接打开查看可视化对比结果）**\n>\n> %s\n\n", webURL))
	}
	if exportURL != "" {
		exportTypeName := "PDF 可视化报告"
		if exportType == 1 {
			exportTypeName = "Excel 差异明细"
		}
		sb.WriteString(fmt.Sprintf("> **📥 %s下载链接**\n>\n> %s\n\n", exportTypeName, exportURL))
	}
	if webURL == "" && exportURL == "" {
		sb.WriteString("> ⚠ 未获取到在线预览或导出链接\n\n")
	}

	// 差异详情
	if len(resp.ComparisonDetail) > 0 {
		sb.WriteString("## 差异详情\n\n")

		// 按类型分组
		var addItems, changeItems, deleteItems []*ess.ComparisonDetail
		for _, d := range resp.ComparisonDetail {
			if d.ComparisonType == nil {
				continue
			}
			switch *d.ComparisonType {
			case "add":
				addItems = append(addItems, d)
			case "change":
				changeItems = append(changeItems, d)
			case "delete":
				deleteItems = append(deleteItems, d)
			}
		}

		if len(changeItems) > 0 {
			sb.WriteString("### 🟡 修改内容\n\n")
			sb.WriteString("| # | 内容类型 | 原文内容 | 新版内容 |\n")
			sb.WriteString("| --- | --- | --- | --- |\n")
			for i, d := range changeItems {
				contentType := getContentTypeName(d.ContentType)
				originText := safeStr(d.OriginText)
				diffText := safeStr(d.DiffText)
				sb.WriteString(fmt.Sprintf("| %d | %s | %s | %s |\n", i+1, contentType, escapeTableCell(originText), escapeTableCell(diffText)))
			}
			sb.WriteString("\n")
		}

		if len(addItems) > 0 {
			sb.WriteString("### 🟢 新增内容\n\n")
			sb.WriteString("| # | 内容类型 | 新增内容 |\n")
			sb.WriteString("| --- | --- | --- |\n")
			for i, d := range addItems {
				contentType := getContentTypeName(d.ContentType)
				diffText := safeStr(d.DiffText)
				if diffText == "" {
					diffText = safeStr(d.OriginText)
				}
				sb.WriteString(fmt.Sprintf("| %d | %s | %s |\n", i+1, contentType, escapeTableCell(diffText)))
			}
			sb.WriteString("\n")
		}

		if len(deleteItems) > 0 {
			sb.WriteString("### 🔴 删除内容\n\n")
			sb.WriteString("| # | 内容类型 | 删除内容 |\n")
			sb.WriteString("| --- | --- | --- |\n")
			for i, d := range deleteItems {
				contentType := getContentTypeName(d.ContentType)
				originText := safeStr(d.OriginText)
				if originText == "" {
					originText = safeStr(d.DiffText)
				}
				sb.WriteString(fmt.Sprintf("| %d | %s | %s |\n", i+1, contentType, escapeTableCell(originText)))
			}
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("## 差异详情\n\n> 未获取到详细差异内容，请查看在线预览页面或下载导出报告。\n\n")
	}

	return sb.String()
}

func getContentTypeName(ct *string) string {
	if ct == nil {
		return "未知"
	}
	switch *ct {
	case "text":
		return "文本"
	case "table":
		return "表格"
	case "picture":
		return "图片"
	default:
		return *ct
	}
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func escapeTableCell(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "|", "\\|")
	if len(s) > 200 {
		s = s[:200] + "..."
	}
	return s
}
