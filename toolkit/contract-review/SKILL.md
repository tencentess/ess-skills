---
name: contract-review
description: |
  腾讯电子签合同审查工作流。当用户要求审查合同、检查合同风险、分析合同条款时使用此 skill。
  支持上传 PDF 自动完成 AI 风险识别，输出风险报告和审查摘要。
version: 1.0.0
---

# 合同审查工作流 Skill

一键完成合同审查全流程: 上传 PDF → 创建审查任务 → 等待 AI 分析 → 获取风险报告。

> **注意**: 合同审查是长耗时任务，AI 分析通常需要 1~10 分钟，请耐心等待，不要中途取消。默认超时为 10 分钟。

## ⚠️ 开发版说明

此 SKILL.md 位于源码目录，仅供开发调试使用。对外发布版请参见 `skills/contract-review/SKILL.md`。

## 前置配置

参见 [凭证配置说明](../references/credentials.md)。

## 使用方法

> **重要**: 所有命令必须在 `contract-review/` 目录下执行（即 `cd toolkit/contract-review`）。

### 一键审查

```bash
bash scripts/run.sh review-workflow \
  --file="/path/to/contract.pdf" \
  --policy-type=0 \
  --export=2
```

**Windows (PowerShell)**:

```powershell
.\scripts\run.ps1 review-workflow `
  --file="C:\path\to\contract.pdf" `
  --policy-type=0 `
  --export=2
```

> **提示**: 若出现"禁止运行脚本"错误，请先执行 `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned`。

### 参数说明

| 参数 | 必填 | 说明 |
|------|------|------|
| `--file` | 是 | 本地 PDF 文件路径 |
| `--policy-type` | 否 | 审查立场: 0=严格 1=中立 2=宽松 |
| `--checklist-id` | 否 | 审查清单 ID，不传则 AI 自动匹配 |
| `--comment` | 否 | 备注 |
| `--export` | 否 | 导出类型: 0=不导出 1=批注文件 2=Excel (默认0) |
| `--timeout` | 否 | 超时时间 (默认10m) |

### 工作流步骤

1. **上传文件** — 将本地 PDF 上传到电子签获取 ResourceId
2. **创建任务** — 调用 CreateBatchContractReviewTask
3. **轮询等待** — 每 5-30 秒检查一次（指数退避），直到 Status=4(成功) 或 Status=5(失败)
4. **获取结果** — 返回风险列表和摘要
5. **导出报告** (可选) — 下载批注文件或 Excel 摘要

### 输出格式

```json
{
  "success": true,
  "data": {
    "task_id": "yD3awUUckpm30d6cUEGhI7Wwu8OHcbXY",
    "review": {
      "Status": 4,
      "TotalRiskCount": 24,
      "HighRiskCount": 13,
      "Risks": [...],
      "Summaries": [...]
    },
    "export": {
      "FileUrl": "https://file.ess.tencent.cn/..."
    }
  },
  "error": null
}
```
