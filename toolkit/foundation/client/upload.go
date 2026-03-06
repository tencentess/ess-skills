package client

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ess "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ess/v20201111"
)

// 允许上传的文件扩展名白名单
var allowedFileExtensions = map[string]bool{
	".pdf": true, ".doc": true, ".docx": true,
	".jpg": true, ".jpeg": true, ".png": true,
	".xls": true, ".xlsx": true, ".html": true,
}

const maxFileSize = 50 * 1024 * 1024 // 50MB

// UploadLocalFile 上传本地文件到电子签，返回 ResourceId（默认 BusinessType 为 DOCUMENT）
func (c *EssClient) UploadLocalFile(filePath string) (string, error) {
	return c.UploadLocalFileWithType(filePath, "DOCUMENT")
}

// UploadLocalFileWithType 上传本地文件到电子签（可指定 BusinessType），返回 ResourceId
func (c *EssClient) UploadLocalFileWithType(filePath, businessType string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" || !allowedFileExtensions[ext] {
		return "", fmt.Errorf("不支持的文件类型 '%s'，支持: %s", ext, getAllowedExtensions())
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("文件访问失败: %w", err)
	}
	if info.Size() > maxFileSize {
		return "", fmt.Errorf("文件大小 (%d bytes) 超过限制 (最大 %dMB)", info.Size(), maxFileSize/1024/1024)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	return c.uploadFileBytes(filepath.Base(filePath), data, businessType)
}

// UploadFileFromURL 从 URL 下载文件并上传到电子签，返回 ResourceId
func (c *EssClient) UploadFileFromURL(fileURL, fileName string) (string, error) {
	httpClient := &http.Client{Timeout: 60 * time.Second}
	resp, err := httpClient.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("下载文件失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载文件 HTTP 状态码异常: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxFileSize+1))
	if err != nil {
		return "", fmt.Errorf("读取下载内容失败: %w", err)
	}
	if len(data) > maxFileSize {
		return "", fmt.Errorf("文件大小超过限制 (最大 %dMB)", maxFileSize/1024/1024)
	}

	if fileName == "" {
		fileName = "uploaded-file.pdf"
	}
	return c.uploadFileBytes(fileName, data, "DOCUMENT")
}

// uploadFileBytes 使用 SDK 上传字节到电子签
func (c *EssClient) uploadFileBytes(fileName string, data []byte, businessType string) (string, error) {
	if businessType == "" {
		businessType = "DOCUMENT"
	}

	uploadClient, err := c.NewUploadClient()
	if err != nil {
		return "", fmt.Errorf("创建上传客户端失败: %w", err)
	}

	fileBody := base64.StdEncoding.EncodeToString(data)

	request := ess.NewUploadFilesRequest()
	request.BusinessType = common.StringPtr(businessType)
	request.Caller = &ess.Caller{
		OperatorId: common.StringPtr(c.Cred.OperatorID),
	}
	request.FileInfos = []*ess.UploadFile{
		{
			FileName: common.StringPtr(fileName),
			FileBody: common.StringPtr(fileBody),
		},
	}

	response, err := uploadClient.UploadFiles(request)
	if err != nil {
		return "", fmt.Errorf("上传文件 API 调用失败: %w", err)
	}

	if response.Response == nil || len(response.Response.FileIds) == 0 {
		return "", fmt.Errorf("上传成功但未返回 FileId")
	}
	return *response.Response.FileIds[0], nil
}

func getAllowedExtensions() string {
	exts := make([]string, 0, len(allowedFileExtensions))
	for ext := range allowedFileExtensions {
		exts = append(exts, ext)
	}
	return strings.Join(exts, ", ")
}
