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
	fs := flag.NewFlagSet("compare-query", flag.ExitOnError)

	var (
		secretID   string
		secretKey  string
		operatorID string
		profile    string
		env        string
		taskID     string
	)

	fs.StringVar(&secretID, "secret-id", "", "腾讯云 SecretId")
	fs.StringVar(&secretKey, "secret-key", "", "腾讯云 SecretKey")
	fs.StringVar(&operatorID, "operator-id", "", "经办人 UserId")
	fs.StringVar(&profile, "profile", "", "配置文件 profile 名称")
	fs.StringVar(&env, "env", "", "环境: test/online (默认 online)")
	fs.StringVar(&taskID, "task-id", "", "对比任务ID (必填)")
	fs.Parse(os.Args[1:])

	if taskID == "" {
		output.PrintError("InvalidParameter", "--task-id 为必填参数")
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

	request := ess.NewDescribeContractComparisonTaskRequest()
	request.Operator = essClient.Operator()
	request.TaskId = common.StringPtr(taskID)

	response, err := essClient.SDKClient.DescribeContractComparisonTask(request)
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
