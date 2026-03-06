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
	fs := flag.NewFlagSet("review-workflow", flag.ExitOnError)

	var (
		secretID    string
		secretKey   string
		operatorID  string
		profile     string
		env         string
		filePath    string
		policyType  int
		checklistID string
		comment     string
		exportType  int
		timeout     string
	)

	fs.StringVar(&secretID, "secret-id", "", "腾讯云 SecretId")
	fs.StringVar(&secretKey, "secret-key", "", "腾讯云 SecretKey")
	fs.StringVar(&operatorID, "operator-id", "", "经办人 UserId")
	fs.StringVar(&profile, "profile", "", "配置文件 profile 名称")
	fs.StringVar(&env, "env", "", "环境: test/online (默认 online)")
	fs.StringVar(&filePath, "file", "", "本地 PDF 文件路径 (必填)")
	fs.IntVar(&policyType, "policy-type", -1, "审查立场: 0-严格 1-中立 2-宽松")
	fs.StringVar(&checklistID, "checklist-id", "", "审查清单ID")
	fs.StringVar(&comment, "comment", "", "备注")
	fs.IntVar(&exportType, "export", 0, "导出类型: 0=不导出 1=批注文件 2=Excel (默认0)")
	fs.StringVar(&timeout, "timeout", "10m", "超时时间 (默认10m)")
	fs.Parse(os.Args[1:])

	if filePath == "" {
		output.PrintError("InvalidParameter", "--file 为必填参数，请提供本地 PDF 文件路径")
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

	// 2. 上传文件
	fmt.Fprintf(os.Stderr, "📤 上传文件: %s\n", filePath)
	resourceID, err := essClient.UploadLocalFile(filePath)
	if err != nil {
		output.PrintErrorf("UploadError", "文件上传失败: %v", err)
	}
	fmt.Fprintf(os.Stderr, "✅ 上传成功, ResourceId: %s\n", resourceID)

	// 3. 创建审查任务
	fmt.Fprintln(os.Stderr, "📝 创建审查任务...")
	createReq := ess.NewCreateBatchContractReviewTaskRequest()
	createReq.Operator = essClient.Operator()
	createReq.ResourceIds = []*string{common.StringPtr(resourceID)}
	if policyType >= 0 {
		createReq.PolicyType = common.Int64Ptr(int64(policyType))
	}
	if checklistID != "" {
		createReq.ChecklistId = common.StringPtr(checklistID)
	}
	if comment != "" {
		createReq.Comment = common.StringPtr(comment)
	}

	createResp, err := essClient.SDKClient.CreateBatchContractReviewTask(createReq)
	if err != nil {
		output.PrintErrorf("CreateTaskError", "创建审查任务失败: %v", err)
	}
	if createResp.Response == nil || len(createResp.Response.TaskIds) == 0 {
		output.PrintError("CreateTaskError", "创建任务成功但未返回 TaskId")
	}
	taskID := *createResp.Response.TaskIds[0]
	fmt.Fprintf(os.Stderr, "✅ 任务创建成功, TaskId: %s\n", taskID)

	// 4. 轮询等待
	fmt.Fprintln(os.Stderr, "⏳ 等待审查完成...")
	pollCfg := poller.DefaultPollConfig
	pollCfg.Timeout = timeoutDuration

	pollResult, err := poller.Poll(context.Background(), pollCfg, func() poller.PollResult {
		queryReq := ess.NewDescribeContractReviewTaskRequest()
		queryReq.Operator = essClient.Operator()
		queryReq.TaskId = common.StringPtr(taskID)

		queryResp, err := essClient.SDKClient.DescribeContractReviewTask(queryReq)
		if err != nil {
			return poller.PollResult{Err: fmt.Errorf("查询任务状态失败: %w", err)}
		}

		if queryResp.Response == nil || queryResp.Response.Status == nil {
			return poller.PollResult{Err: fmt.Errorf("查询任务返回空响应")}
		}

		status := *queryResp.Response.Status
		switch status {
		case 4: // 成功
			return poller.PollResult{Done: true, Data: queryResp}
		case 5: // 失败
			return poller.PollResult{Err: fmt.Errorf("审查任务执行失败")}
		default:
			return poller.PollResult{Done: false}
		}
	})
	if err != nil {
		output.PrintErrorf("PollError", "等待审查完成失败: %v", err)
	}

	fmt.Fprintln(os.Stderr, "✅ 审查完成!")

	reviewResp := pollResult.(*ess.DescribeContractReviewTaskResponse)

	// 5. 导出带批注文件（始终导出）
	fmt.Fprintln(os.Stderr, "📥 导出带批注文件...")
	var annotatedFileURL string
	exportAnnotReq := ess.NewExportContractReviewResultRequest()
	exportAnnotReq.Operator = essClient.Operator()
	exportAnnotReq.TaskId = common.StringPtr(taskID)
	exportAnnotReq.FileType = common.Int64Ptr(1) // 1=带风险批注文件
	exportAnnotResp, err := essClient.SDKClient.ExportContractReviewResult(exportAnnotReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠ 导出带批注文件失败: %v\n", err)
	} else if exportAnnotResp.Response != nil && exportAnnotResp.Response.Url != nil {
		annotatedFileURL = *exportAnnotResp.Response.Url
		fmt.Fprintln(os.Stderr, "✅ 带批注文件导出成功!")
	}

	// 6. 可选额外导出（用户指定 exportType=2 时导出 Excel）
	var excelFileURL string
	if exportType == 2 {
		fmt.Fprintln(os.Stderr, "📥 导出审查结果 Excel...")
		exportExcelReq := ess.NewExportContractReviewResultRequest()
		exportExcelReq.Operator = essClient.Operator()
		exportExcelReq.TaskId = common.StringPtr(taskID)
		exportExcelReq.FileType = common.Int64Ptr(2) // 2=审查结果&摘要(xlsx)
		exportExcelResp, err := essClient.SDKClient.ExportContractReviewResult(exportExcelReq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠ 导出 Excel 失败: %v\n", err)
		} else if exportExcelResp.Response != nil && exportExcelResp.Response.Url != nil {
			excelFileURL = *exportExcelResp.Response.Url
			fmt.Fprintln(os.Stderr, "✅ Excel 导出成功!")
		}
	}

	// 7. 获取在线审查页面链接
	fmt.Fprintln(os.Stderr, "🔗 获取在线审查页面链接...")
	var webURL string
	webUrlReq := ess.NewDescribeContractReviewWebUrlRequest()
	webUrlReq.Operator = essClient.Operator()
	webUrlReq.TaskId = common.StringPtr(taskID)
	webUrlResp, err := essClient.SDKClient.DescribeContractReviewWebUrl(webUrlReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠ 获取在线审查页面链接失败: %v\n", err)
	} else if webUrlResp.Response != nil && webUrlResp.Response.WebUrl != nil {
		webURL = *webUrlResp.Response.WebUrl
		fmt.Fprintln(os.Stderr, "✅ 在线审查页面链接获取成功!")
	}

	// 8. 构建 Markdown 结果
	md := renderMarkdown(reviewResp, taskID, filepath.Base(filePath), policyType, annotatedFileURL, excelFileURL, webURL)

	// 直接将 Markdown 输出到 stderr，让终端/Agent 可以直接展示
	fmt.Fprintln(os.Stderr, "\n"+md)

	// 构建结构化输出（stdout JSON 供程序消费）
	finalResult := map[string]interface{}{
		"task_id":  taskID,
		"markdown": md,
	}
	if annotatedFileURL != "" {
		finalResult["annotated_file_url"] = annotatedFileURL
	}
	if excelFileURL != "" {
		finalResult["excel_file_url"] = excelFileURL
	}
	if webURL != "" {
		finalResult["web_url"] = webURL
	}

	output.PrintSuccess(finalResult)
}

// renderMarkdown 将审查结果渲染为 Markdown 格式
func renderMarkdown(resp *ess.DescribeContractReviewTaskResponse, taskID, fileName string, policyType int, annotatedURL, excelURL, webURL string) string {
	var sb strings.Builder

	sb.WriteString("# 合同智能审查报告\n\n")

	// 基本信息
	sb.WriteString("## 基本信息\n\n")
	sb.WriteString(fmt.Sprintf("| 项目 | 内容 |\n"))
	sb.WriteString(fmt.Sprintf("| --- | --- |\n"))
	sb.WriteString(fmt.Sprintf("| 文件名称 | %s |\n", fileName))
	sb.WriteString(fmt.Sprintf("| 任务ID | `%s` |\n", taskID))
	policyName := map[int]string{0: "严格", 1: "中立", 2: "宽松"}[policyType]
	if policyName == "" {
		policyName = "默认"
	}
	sb.WriteString(fmt.Sprintf("| 审查立场 | %s |\n", policyName))

	r := resp.Response

	// 风险统计
	if r.TotalRiskCount != nil {
		sb.WriteString(fmt.Sprintf("| 风险总数 | **%d** |\n", *r.TotalRiskCount))
	}
	if r.HighRiskCount != nil {
		sb.WriteString(fmt.Sprintf("| 高风险数 | **%d** |\n", *r.HighRiskCount))
	}
	sb.WriteString("\n")

	// 重要链接（重点显示）
	sb.WriteString("## 📌 重要链接\n\n")
	if webURL != "" {
		sb.WriteString(fmt.Sprintf("> **🌐 在线审查页面（可直接打开查看完整审查结果）**\n>\n> %s\n\n", webURL))
	} else {
		sb.WriteString("> ⚠ 未获取到在线审查页面链接\n\n")
	}
	if annotatedURL != "" {
		sb.WriteString(fmt.Sprintf("> **📄 带风险批注的文档下载链接**\n>\n> %s\n\n", annotatedURL))
	} else {
		sb.WriteString("> ⚠ 未获取到带批注文档下载链接\n\n")
	}
	if excelURL != "" {
		sb.WriteString(fmt.Sprintf("> **📊 审查结果 Excel 下载链接**\n>\n> %s\n\n", excelURL))
	}

	// 合同摘要
	if len(r.Summaries) > 0 {
		sb.WriteString("## 合同摘要\n\n")
		for _, summary := range r.Summaries {
			if summary.Name == nil {
				continue
			}
			categoryName := map[string]string{
				"Base":        "合同信息",
				"Identity":    "主体信息",
				"Performance": "履约条款",
			}[*summary.Name]
			if categoryName == "" {
				categoryName = *summary.Name
			}
			sb.WriteString(fmt.Sprintf("### %s\n\n", categoryName))
			sb.WriteString("| 字段 | 内容 |\n")
			sb.WriteString("| --- | --- |\n")
			for _, info := range summary.Infos {
				key := ""
				value := ""
				if info.Key != nil {
					key = *info.Key
				}
				if info.Value != nil {
					value = *info.Value
				}
				sb.WriteString(fmt.Sprintf("| %s | %s |\n", key, value))
			}
			sb.WriteString("\n")
		}
	}

	// 风险详情
	if len(r.Risks) > 0 {
		sb.WriteString("## 风险详情\n\n")

		// 按风险等级分组
		var highRisks, normalRisks []*ess.OutputRisk
		for _, risk := range r.Risks {
			if risk.RiskLevel != nil && *risk.RiskLevel == "HIGH" {
				highRisks = append(highRisks, risk)
			} else {
				normalRisks = append(normalRisks, risk)
			}
		}

		if len(highRisks) > 0 {
			sb.WriteString("### 🔴 高风险\n\n")
			for i, risk := range highRisks {
				writeRiskItem(&sb, i+1, risk)
			}
		}

		if len(normalRisks) > 0 {
			sb.WriteString("### 🟡 一般风险\n\n")
			for i, risk := range normalRisks {
				writeRiskItem(&sb, i+1, risk)
			}
		}
	} else {
		sb.WriteString("## 风险详情\n\n✅ 未发现风险项\n\n")
	}

	return sb.String()
}

// writeRiskItem 写入单个风险项
func writeRiskItem(sb *strings.Builder, idx int, risk *ess.OutputRisk) {
	name := ""
	if risk.RiskName != nil {
		name = *risk.RiskName
	}
	sb.WriteString(fmt.Sprintf("**%d. %s**\n\n", idx, name))

	if risk.RiskDescription != nil && *risk.RiskDescription != "" {
		sb.WriteString(fmt.Sprintf("- **风险描述**: %s\n", *risk.RiskDescription))
	}
	if risk.Content != nil && *risk.Content != "" {
		sb.WriteString(fmt.Sprintf("- **原文内容**: %s\n", *risk.Content))
	}
	if risk.RiskAdvice != nil && *risk.RiskAdvice != "" {
		sb.WriteString(fmt.Sprintf("- **修改建议**: %s\n", *risk.RiskAdvice))
	}
	if risk.RiskBasis != nil && *risk.RiskBasis != "" {
		sb.WriteString(fmt.Sprintf("- **审查依据**: %s\n", *risk.RiskBasis))
	}
	if len(risk.RiskPresentation) > 0 {
		sb.WriteString("- **风险评估**:\n")
		for _, p := range risk.RiskPresentation {
			if p != nil {
				sb.WriteString(fmt.Sprintf("  - %s\n", *p))
			}
		}
	}
	sb.WriteString("\n")
}
