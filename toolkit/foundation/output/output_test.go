package output

import (
	"encoding/json"
	"testing"
)

func TestPrintSuccess_Format(t *testing.T) {
	data := map[string]string{"task_id": "abc123"}
	dataBytes, _ := json.Marshal(data)

	resp := Response{
		Success: true,
		Data:    dataBytes,
		Error:   nil,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var parsed Response
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !parsed.Success {
		t.Error("expected success=true")
	}
	if parsed.Error != nil {
		t.Error("expected error=nil")
	}

	var innerData map[string]string
	if err := json.Unmarshal(parsed.Data, &innerData); err != nil {
		t.Fatalf("unmarshal data failed: %v", err)
	}
	if innerData["task_id"] != "abc123" {
		t.Errorf("expected task_id=abc123, got %s", innerData["task_id"])
	}
}

func TestErrorFormat(t *testing.T) {
	resp := Response{
		Success: false,
		Data:    json.RawMessage("null"),
		Error:   &ErrorInfo{Code: "AuthFailure", Message: "invalid key"},
	}

	out, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var parsed Response
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if parsed.Success {
		t.Error("expected success=false")
	}
	if parsed.Error == nil {
		t.Fatal("expected error != nil")
	}
	if parsed.Error.Code != "AuthFailure" {
		t.Errorf("expected code=AuthFailure, got %s", parsed.Error.Code)
	}
}
