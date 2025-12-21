package main

import "testing"

func TestGenerateFieldListSignatureWithTimestamp_KnownSample(t *testing.T) {
	got, err := GenerateFieldListSignatureWithTimestamp(
		"20251219",
		"2025082802482655",
		"5003000103",
		"1002",
		"o_9oO5UjYM1frKP537iCEGv0JID4",
		1766054837772,
	)
	if err != nil {
		t.Fatalf("GenerateFieldListSignatureWithTimestamp() error = %v", err)
	}

	want := "apiKey=e98ce2565b09ecc0&timestamp=1766054837772&channelId=11&netUserId=2025082802482655&venueId=5003000103&serviceId=1002&day=20251219&selectByfullTag=0&centerId=50030001&fieldType=1837&tenantId=82&openId=o_9oO5UjYM1frKP537iCEGv0JID4&version=1&sign=6030ee1244875efa2de392ae62924da1"
	if got != want {
		t.Fatalf("signature mismatch:\n got: %s\nwant: %s", got, want)
	}
}
