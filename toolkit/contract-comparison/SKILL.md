---
name: contract-comparison
description: |
  腾讯电子签合同对比工作流。当用户要求对比两份合同、比较合同差异、查看合同修改内容时使用此 skill。
  支持上传两份文件（PDF、Word 等格式）自动完成差异分析，输出新增、修改、删除的差异明细。
version: 1.0.0
---

# 合同对比工作流 Skill

一键完成合同对比全流程: 上传两份文件（支持 PDF、Word 等格式）→ 创建对比任务 → 等待分析 → 获取差异报告。

> **注意**: 合同对比是长耗时任务，分析通常需要 1~10 分钟，请耐心等待，不要中途取消。默认超时为 10 分钟。

## ⚠️ 开发版说明

此 SKILL.md 位于源码目录，仅供开发调试使用。对外发布版请参见 `skills/contract-comparison/SKILL.md`。

## 前置配置

参见 [凭证配置说明](../references/credentials.md)。

## 使用方法

> **重要**: 所有命令必须在 `contract-comparison/` 目录下执行（即 `cd toolkit/contract-comparison`）。

### 一键对比

```bash
bash scripts/run.sh compare-workflow \
  --origin-file="/path/to/old-contract.pdf" \
  --diff-file="/path/to/new-contract.docx" \
  --export=0
```

**Windows (PowerShell)**:

```powershell
.\scripts\run.ps1 compare-workflow `
  --origin-file="C:\path\to\old-contract.pdf" `
  --diff-file="C:\path\to\new-contract.docx" `
  --export=0
```

> **提示**: 若出现"禁止运行脚本"错误，请先执行 `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned`。

### 参数说明

| 参数 | 必填 | 说明 |
|------|------|------|
| `--origin-file` | 是 | 原版文件路径（支持 PDF、Word 等格式，≤60M） |
| `--diff-file` | 是 | 新版文件路径（支持 PDF、Word 等格式，≤60M） |
| `--comment` | 否 | 备注 |
| `--export` | 否 | 导出: -1=不导出 0=PDF报告 1=Excel明细 (默认-1) |
| `--timeout` | 否 | 超时时间 (默认10m) |

### 工作流步骤

1. **上传两份文件** — 分别上传原版和新版文件（PDF、Word 等）
2. **创建对比任务** — 调用 CreateContractComparisonTask
3. **轮询等待** — 每 5-30 秒检查一次（指数退避），直到 Status=2(成功) 或 Status=3(失败)
4. **获取结果** — 返回差异统计和详情
5. **导出报告** (可选) — 下载 PDF 可视化报告或 Excel 差异明细

### 输出格式

```json
{
  "success": true,
  "data": {
    "task_id": "yDtTFUUckp9elogcUudRpd6uv7cdx6Qa",
    "comparison": {
      "Status": 2,
      "TotalDiffCount": 12,
      "AddDiffCount": 7,
      "ChangeDiffCount": 2,
      "DeleteDiffCount": 3
    },
    "export": {
      "FileUrl": "https://file.ess.tencent.cn/..."
    }
  },
  "error": null
}
```
