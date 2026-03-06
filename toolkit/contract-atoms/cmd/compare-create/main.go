package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/tencentess/ess-skills/toolkit/foundation/client"
	"github.com/tencentess/ess-skills/toolkit/foundation/output"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ess "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ess/v20201111"
)

func main() {
	fs := flag.NewFlagSet("compare-create", flag.ExitOnError)

	var (
		secretID     string
		secretKey    string
		operatorID   string
		profile      string
		env          string
		originFileID string
		diffFileID   string
		comment      string
	)

	fs.StringVar(&secretID, "secret-id", "", "腾讯云 SecretId")
	fs.StringVar(&secretKey, "secret-key", "", "腾讯云 SecretKey")
	fs.StringVar(&operatorID, "operator-id", "", "经办人 UserId")
	fs.StringVar(&profile, "profile", "", "配置文件 profile 名称")
	fs.StringVar(&env, "env", "", "环境: test/online (默认 online)")
	fs.StringVar(&originFileID, "origin-file-id", "", "原版文件资源ID (必填)")
	fs.StringVar(&diffFileID, "diff-file-id", "", "新版文件资源ID (必填)")
	fs.StringVar(&comment, "comment", "", "备注")
	fs.Parse(os.Args[1:])

	if originFileID == "" || diffFileID == "" {
		output.PrintError("InvalidParameter", "--origin-file-id 和 --diff-file-id 为必填参数")
	}

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

	request := ess.NewCreateContractComparisonTaskRequest()
	request.Operator = essClient.Operator()
	request.OriginFileResourceId = common.StringPtr(originFileID)
	request.DiffFileResourceId = common.StringPtr(diffFileID)
	if comment != "" {
		request.Comment = common.StringPtr(comment)
	}

	response, err := essClient.SDKClient.CreateContractComparisonTask(request)
	if err != nil {
		output.PrintError("ApiError", err.Error())
	}

	var resp map[string]interface{}
	json.Unmarshal([]byte(response.ToJsonString()), &resp)
	if r, ok := resp["Response"]; ok {
		output.PrintSuccess(r)
	} else {
		output.PrintSuccess(resp)
	}
}
