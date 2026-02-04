package main

import (
	"testing"
)

func TestBuildPayOrderBodySignature(t *testing.T) {
	// 设置测试数据（只设置变量,常量已在主文件定义）
	apiSecret = "UzpnKLpYdQJgVBxQFeKyxw=="
	apiVersion = 31
	openId = "o_9oO5YDtVXIBdWsrliOwuBEXi3M"

	// 测试参数
	timestamp := int64(1770163261959)
	tradeId := "2026020493646400"
	userId := "2025101002519673"
	userEcardNo := "E5003000100037657"

	// 期望的签名
	expectedSign := "9c0cd3ae4fd9d517577ad651a223e354"

	// 生成签名字符串
	body := buildPayOrderBody(tradeId, timestamp, userId, userEcardNo)

	t.Logf("生成的请求体: %s", body)

	// 从生成的body中提取sign参数
	// 期望格式: apiKey=...&timestamp=...&sign=xxx
	var actualSign string
	for _, part := range splitParams(body) {
		if len(part) > 5 && part[:5] == "sign=" {
			actualSign = part[5:]
			break
		}
	}

	t.Logf("期望签名: %s", expectedSign)
	t.Logf("实际签名: %s", actualSign)

	if actualSign != expectedSign {
		t.Errorf("签名不匹配!\n期望: %s\n实际: %s", expectedSign, actualSign)
	} else {
		t.Log("✅ 签名验证成功!")
	}
}

// TestGenerateSignDetail 测试签名生成的详细过程
func TestGenerateSignDetail(t *testing.T) {
	// 设置测试数据
	apiSecret = "UzpnKLpYdQJgVBxQFeKyxw=="
	apiVersion = 31
	openId = "o_9oO5YDtVXIBdWsrliOwuBEXi3M"

	// 测试参数
	timestamp := int64(1770163261959)
	tradeId := "2026020493646400"
	userId := "2025101002519673"
	userEcardNo := "E5003000100037657"

	// 构建参数
	params := map[string]string{
		"netUserId": userId,
		"tradeId":   tradeId,
		"payGroup":  "[]",
		"ecardNo":   userEcardNo,
		"openId":    openId,
	}

	// 生成签名
	sign := generateSign("/aisports-api/api/pay/payOrder", params, timestamp)

	t.Logf("API路径: /aisports-api/api/pay/payOrder")
	t.Logf("参数:")
	t.Logf("  apiKey: %s", APIKey)
	t.Logf("  timestamp: %d", timestamp)
	t.Logf("  channelId: %s", ChannelID)
	t.Logf("  netUserId: %s", userId)
	t.Logf("  tradeId: %s", tradeId)
	t.Logf("  payGroup: []")
	t.Logf("  ecardNo: %s", userEcardNo)
	t.Logf("  centerId: %s", CenterID)
	t.Logf("  tenantId: %s", TenantID)
	t.Logf("  openId: %s", openId)
	t.Logf("  version: %d", apiVersion)
	t.Logf("apiSecret: %s", apiSecret)
	t.Logf("生成的签名: %s", sign)

	expectedSign := "9c0cd3ae4fd9d517577ad651a223e354"
	if sign != expectedSign {
		t.Errorf("签名不匹配!\n期望: %s\n实际: %s", expectedSign, sign)
	} else {
		t.Log("✅ 签名验证成功!")
	}
}

// 辅助函数：分割参数字符串
func splitParams(s string) []string {
	var result []string
	var current string
	for i := 0; i < len(s); i++ {
		if s[i] == '&' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(s[i])
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
