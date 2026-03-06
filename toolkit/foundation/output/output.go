package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// Response 统一输出结构
type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *ErrorInfo      `json:"error"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// PrintSuccess 输出成功结果到 stdout
func PrintSuccess(data interface{}) {
	dataBytes, _ := json.Marshal(data)
	resp := Response{
		Success: true,
		Data:    dataBytes,
		Error:   nil,
	}
	out, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintln(os.Stdout, string(out))
}

// PrintError 输出错误结果到 stdout 并以退出码 1 退出
func PrintError(code, message string) {
	resp := Response{
		Success: false,
		Data:    json.RawMessage("null"),
		Error:   &ErrorInfo{Code: code, Message: message},
	}
	out, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintln(os.Stdout, string(out))
	os.Exit(1)
}

// PrintErrorf 格式化错误信息输出
func PrintErrorf(code, format string, args ...interface{}) {
	PrintError(code, fmt.Sprintf(format, args...))
}
