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
	)
	if err != nil {
		t.Fatalf("GenerateFieldListSignatureWithTimestamp() error = %v", err)
	}

	want := "apiKey=e98ce2565b09ecc0&timestamp=1767100063371&channelId=11&netUserId=2025082802482655&venueId=5003000103&serviceId=1002&day=20251231&selectByfullTag=0&centerId=50030001&fieldType=1837&tenantId=82&openId=o_9oO5UjYM1frKP537iCEGv0JID4&version=9&sign=fedeb01a1e711b1785e51b8466604448"
	if got != want {
		t.Fatalf("signature mismatch:\n got: %s\nwant: %s", got, want)
	}
}
