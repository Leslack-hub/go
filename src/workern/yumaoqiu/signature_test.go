package main

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

// TestSignatureConsistency 测试Go版本与JavaScript版本的签名一致性
func TestSignatureConsistency(t *testing.T) {
	// 使用固定时间戳确保两个版本使用相同的输入
	fixedTimestamp := "1757076700000"

	testCases := []struct {
		name      string
		method    string
		day       string
		fieldInfo string
		netUserID string
		venueID   string
		serviceID string
	}{
		{
			name:      "fieldList签名测试",
			method:    "fieldList",
			day:       "20250830",
			netUserID: "2025082802482655",
			venueID:   "5003000101",
			serviceID: "1002",
		},
		{
			name:      "newOrder签名测试",
			method:    "newOrder",
			day:       "20250830",
			fieldInfo: "f3a74ed8bf6efbe143e33c0ba1cf9e26,fade00ab4725896da1840e9c0125dc0f",
			netUserID: "2025082802482655",
			venueID:   "5003000101",
			serviceID: "1002",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 生成Go版本的签名（使用固定时间戳）
			var goSignature string
			var err error

			// 解析固定时间戳
			timestampInt, parseErr := strconv.ParseInt(fixedTimestamp, 10, 64)
			if parseErr != nil {
				t.Fatalf("时间戳解析失败: %v", parseErr)
			}

			switch tc.method {
			case "fieldList":
				goSignature, err = GenerateFieldListSignatureWithTimestamp(tc.day, tc.netUserID, tc.venueID, tc.serviceID, timestampInt)
			case "newOrder":
				goSignature, err = GenerateNewOrderSignatureWithTimestamp(tc.day, tc.fieldInfo, tc.netUserID, tc.serviceID, tc.venueID, timestampInt)
			default:
				t.Fatalf("未知的方法: %s", tc.method)
			}

			if err != nil {
				t.Fatalf("Go签名生成失败: %v", err)
			}

			// 生成JavaScript版本的签名（如果Node.js可用）
			var jsCmd *exec.Cmd
			switch tc.method {
			case "fieldList":
				jsCmd = exec.Command("node", "signature_generator.js", "-m", "fieldList", "--day", tc.day, "--netUserId", tc.netUserID, "--venueId", tc.venueID, "--serviceId", tc.serviceID, "--timestamp", fixedTimestamp)
			case "newOrder":
				jsCmd = exec.Command("node", "signature_generator.js", "-m", "newOrder", "--day", tc.day, "--fieldInfo", tc.fieldInfo, "--netUserId", tc.netUserID, "--venueId", tc.venueID, "--serviceId", tc.serviceID, "--timestamp", fixedTimestamp)
			}

			jsOutput, jsErr := jsCmd.Output()
			if jsErr != nil {
				t.Logf("JavaScript版本不可用，跳过对比测试: %v", jsErr)
				t.Logf("Go版本签名结果: %s", goSignature)
				return
			}

			jsSignature := strings.TrimSpace(string(jsOutput))

			// 比较结果
			if goSignature != jsSignature {
				t.Errorf("签名不一致:\nGo版本:  %s\nJS版本:  %s", goSignature, jsSignature)
			} else {
				t.Logf("签名一致: %s", goSignature)
			}
		})
	}
}

// TestSignatureGeneration 测试签名生成功能
func TestSignatureGeneration(t *testing.T) {
	testCases := []struct {
		name      string
		apiPath   string
		params    map[string]interface{}
		expectErr bool
	}{
		{
			name:    "基本签名生成",
			apiPath: "/test/api",
			params: map[string]interface{}{
				"param1": "value1",
				"param2": "value2",
			},
			expectErr: false,
		},
		{
			name:      "空参数签名生成",
			apiPath:   "/test/api",
			params:    map[string]interface{}{},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := generateSignature(tc.apiPath, tc.params, nil)
			if tc.expectErr {
				if err == nil {
					t.Error("期望出现错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Fatalf("签名生成失败: %v", err)
			}

			if result.Sign == "" {
				t.Error("签名为空")
			}

			if result.APIKey != APIKey {
				t.Errorf("API Key不匹配: 期望 %s, 实际 %s", APIKey, result.APIKey)
			}

			if result.ChannelID != ChannelID {
				t.Errorf("Channel ID不匹配: 期望 %s, 实际 %s", ChannelID, result.ChannelID)
			}

			// 测试URL参数转换
			urlParams := toURLParams(result)
			if urlParams == "" {
				t.Error("URL参数为空")
			}

			t.Logf("签名结果: %s", result.Sign)
			t.Logf("URL参数: %s", urlParams)
		})
	}
}

// BenchmarkSignatureGeneration 签名生成性能测试
func BenchmarkSignatureGeneration(b *testing.B) {
	apiPath := "/test/api"
	params := map[string]interface{}{
		"param1": "value1",
		"param2": "value2",
		"param3": "value3",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generateSignature(apiPath, params, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}
