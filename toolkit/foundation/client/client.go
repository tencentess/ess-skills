package client

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ess "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ess/v20201111"
)

// EndpointConfig 存储不同环境的 Endpoint 映射
var endpointConfig = map[string]map[string]string{
	"test": {
		"_default":    "ess.test.ess.tencent.cn",
		"UploadFiles": "file.test.ess.tencent.cn",
	},
	"online": {
		"_default":    "ess.tencentcloudapi.com",
		"UploadFiles": "file.ess.tencent.cn",
	},
}

// EssClient 腾讯电子签 API 客户端（基于官方 SDK）
type EssClient struct {
	Cred      *ResolvedCredentials
	SDKClient *ess.Client
}

// NewEssClient 从已解析的凭证创建 ESS 客户端
func NewEssClient(cred *ResolvedCredentials) (*EssClient, error) {
	if cred == nil {
		return nil, fmt.Errorf("凭证不能为空")
	}

	credential := common.NewCredential(cred.SecretID, cred.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = getEndpoint(cred.Env, "_default")

	sdkClient, err := ess.NewClient(credential, "", cpf)
	if err != nil {
		return nil, fmt.Errorf("创建 SDK 客户端失败: %w", err)
	}

	return &EssClient{
		Cred:      cred,
		SDKClient: sdkClient,
	}, nil
}

// NewUploadClient 创建专用于文件上传的 SDK 客户端（使用 file.ess.tencent.cn 域名）
func (c *EssClient) NewUploadClient() (*ess.Client, error) {
	credential := common.NewCredential(c.Cred.SecretID, c.Cred.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = getEndpoint(c.Cred.Env, "UploadFiles")

	return ess.NewClient(credential, "", cpf)
}

// Operator 返回 SDK 的 UserInfo 结构，自动注入 OperatorID
func (c *EssClient) Operator() *ess.UserInfo {
	return &ess.UserInfo{
		UserId: common.StringPtr(c.Cred.OperatorID),
	}
}

// getEndpoint 根据环境和接口名获取 Endpoint
func getEndpoint(env, action string) string {
	if env == "" {
		env = "online"
	}
	epMap, ok := endpointConfig[env]
	if !ok {
		epMap = endpointConfig["online"]
	}
	if ep, exists := epMap[action]; exists {
		return ep
	}
	return epMap["_default"]
}
