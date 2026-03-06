package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/tencentess/ess-skills/toolkit/foundation/client"
	"github.com/tencentess/ess-skills/toolkit/foundation/output"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ess "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ess/v20201111"
)

func main() {
	fs := flag.NewFlagSet("review-create", flag.ExitOnError)

	var (
		secretID    string
		secretKey   string
		operatorID  string
		profile     string
		env         string
		resourceIDs string
		policyType  int
		checklistID string
		comment     string
	)

	fs.StringVar(&secretID, "secret-id", "", "腾讯云 SecretId")
	fs.StringVar(&secretKey, "secret-key", "", "腾讯云 SecretKey")
	fs.StringVar(&operatorID, "operator-id", "", "经办人 UserId")
	fs.StringVar(&profile, "profile", "", "配置文件 profile 名称")
	fs.StringVar(&env, "env", "", "环境: test/online (默认 online)")
	fs.StringVar(&resourceIDs, "resource-ids", "", "PDF资源ID列表，逗号分隔 (必填)")
	fs.IntVar(&policyType, "policy-type", -1, "审查立场: 0-严格 1-中立 2-宽松 (不传则AI推荐)")
	fs.StringVar(&checklistID, "checklist-id", "", "审查清单ID (不传则AI自动匹配)")
	fs.StringVar(&comment, "comment", "", "备注")
	fs.Parse(os.Args[1:])

	if resourceIDs == "" {
		output.PrintError("InvalidParameter", "--resource-ids 为必填参数")
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

	ids := strings.Split(resourceIDs, ",")
	resIds := make([]*string, len(ids))
	for i := range ids {
		resIds[i] = common.StringPtr(ids[i])
	}

	request := ess.NewCreateBatchContractReviewTaskRequest()
	request.Operator = essClient.Operator()
	request.ResourceIds = resIds
	if policyType >= 0 {
		request.PolicyType = common.Int64Ptr(int64(policyType))
	}
	if checklistID != "" {
		request.ChecklistId = common.StringPtr(checklistID)
	}
	if comment != "" {
		request.Comment = common.StringPtr(comment)
	}

	response, err := essClient.SDKClient.CreateBatchContractReviewTask(request)
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
