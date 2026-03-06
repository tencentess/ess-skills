---
name: contract-atoms
description: |
  腾讯电子签合同智能原子操作。当用户需要单独执行某个合同 API 步骤时使用此 skill，
  如：创建审查/对比任务、查询任务状态、导出结果、查看审查清单。适合需要精细控制流程的场景。
version: 1.0.0
---

# 合同智能原子操作 Skill

提供腾讯电子签合同智能 API 的单步调用能力。每个命令对应一个独立的 API 操作。

## ⚠️ 开发版说明

此 SKILL.md 位于源码目录，仅供开发调试使用。对外发布版请参见 `skills/contract-atoms/SKILL.md`。

## 前置配置

参见 [凭证配置说明](../references/credentials.md)。

## 使用方法

> **重要**: 所有命令必须在 `contract-atoms/` 目录下执行（即 `cd toolkit/contract-atoms`）。

所有命令通过 `scripts/run.sh` 统一入口调用，自动检测当前平台架构：

> **Windows 用户**: 使用 `.\scripts\run.ps1 <command> [args...]` 代替 `bash scripts/run.sh`。
> 若提示"禁止运行脚本"，请先执行 `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned`。

### 审查相关

#### review-create — 创建合同审查任务
```bash
bash scripts/run.sh review-create \
  --resource-ids="yD3a1UUckpm1f4glUxASQ9WvM1liHmXI" \
  --policy-type=0 \
  --checklist-id="可选" \
  --comment="可选备注"
```
参数:
- `--resource-ids` (必填): PDF 资源 ID，多个逗号分隔（最多5个）
- `--policy-type` (可选): 0=严格, 1=中立, 2=宽松
- `--checklist-id` (可选): 审查清单 ID，不传则 AI 自动匹配
- `--comment` (可选): 备注

#### review-query — 查询审查任务结果
```bash
bash scripts/run.sh review-query --task-id="yD3awUUckpm30d6cUEGhI7Wwu8OHcbXY"
```
参数: `--task-id` (必填)
返回: 任务状态 + 风险列表 + 摘要。Status: 1=创建成功, 2=排队中, 3=执行中, 4=成功, 5=失败

#### review-export — 导出审查结果
```bash
bash scripts/run.sh review-export --task-id="xxx" --file-type=2
```
参数:
- `--task-id` (必填): 审查任务ID
- `--file-type` (可选): 1=带风险批注文件, 2=审查结果&摘要xlsx (默认2)

#### review-checklists — 查看企业审查清单
```bash
bash scripts/run.sh review-checklists
```
无业务参数，返回企业已配置的审查清单列表。

### 对比相关

#### compare-create — 创建合同对比任务
```bash
bash scripts/run.sh compare-create \
  --origin-file-id="原版文件资源ID" \
  --diff-file-id="新版文件资源ID" \
  --comment="可选"
```
参数:
- `--origin-file-id` (必填): 原版文件资源ID
- `--diff-file-id` (必填): 新版文件资源ID
- `--comment` (可选): 备注

#### compare-query — 查询对比结果
```bash
bash scripts/run.sh compare-query --task-id="xxx"
```
参数: `--task-id` (必填)
返回: Status: 0=待创建, 1=对比中, 2=成功, 3=失败

#### compare-export — 导出对比报告
```bash
bash scripts/run.sh compare-export --task-id="xxx" --export-type=0
```
参数:
- `--task-id` (必填): 对比任务ID
- `--export-type` (可选): 0=PDF可视化报告, 1=Excel差异明细 (默认0)

## 输出格式

所有命令统一输出 JSON 到 stdout:
```json
{"success": true, "data": {...}, "error": null}
{"success": false, "data": null, "error": {"code": "xxx", "message": "xxx"}}
```
