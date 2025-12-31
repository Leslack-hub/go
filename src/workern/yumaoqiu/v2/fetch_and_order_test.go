package main

import "testing"

func TestGenerateFieldListSignatureWithTimestamp_KnownSample(t *testing.T) {
	got, err := GenerateFieldListSignatureWithTimestamp(
		"20251231",
		"2025082802482655",
		"5003000103",
		"1002",
		"o_9oO5UjYM1frKP537iCEGv0JID4",
		"0gDFZNGBdobPSQjIUbp/NA==",
		9,
		1767100063371,
		"1837",
	)
	if err != nil {
		t.Fatalf("GenerateFieldListSignatureWithTimestamp() error = %v", err)
	}

	want := "apiKey=e98ce2565b09ecc0&timestamp=1767100063371&channelId=11&netUserId=2025082802482655&venueId=5003000103&serviceId=1002&day=20251231&selectByfullTag=0&centerId=50030001&fieldType=1837&tenantId=82&openId=o_9oO5UjYM1frKP537iCEGv0JID4&version=9&sign=fedeb01a1e711b1785e51b8466604448"
	if got != want {
		t.Fatalf("signature mismatch:\n got: %s\nwant: %s", got, want)
	}
}

func TestGenerateNewOrderSignatureWithTimestamp_KnownSample(t *testing.T) {
	result, err := generateSignatureWithTimestamp(
		"/aisports-api/wechatAPI/order/newOrder",
		map[string]any{
			"serviceId": "1002",
			"day":       "20251231",
			"fieldType": "1837",
			"fieldInfo": "cc2c804dd17376217dfcda45ffa08f05",
			"ticket":    "",
			"randStr":   "",
			"venueId":   "5003000103",
			"netUserId": "2025082802482655",
			"openId":    "o_9oO5UjYM1frKP537iCEGv0JID4",
		},
		"0gDFZNGBdobPSQjIUbp/NA==",
		14,
		1767158725485,
	)
	if err != nil {
		t.Fatalf("generateSignatureWithTimestamp() error = %v", err)
	}

	got := toURLParams(result)
	want := "apiKey=e98ce2565b09ecc0&timestamp=1767158725485&channelId=11&venueId=5003000103&serviceId=1002&centerId=50030001&day=20251231&fieldType=1837&fieldInfo=cc2c804dd17376217dfcda45ffa08f05&ticket=&randStr=&netUserId=2025082802482655&tenantId=82&openId=o_9oO5UjYM1frKP537iCEGv0JID4&version=14&sign=f492ba8d6b0bc7ae5b6f1ee4a6571c2e"
	if got != want {
		t.Fatalf("signature mismatch:\n got: %s\nwant: %s", got, want)
	}
}
