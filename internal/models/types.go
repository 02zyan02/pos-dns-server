package models

import "time"

type Binding struct {
	Company   string    `json:"company"`
	Passkey   string    `json:"passkey"`
	ClientIP  string    `json:"client_ip"`
	OTP       string    `json:"otp,omitempty"`
	OTPExpiry time.Time `json:"otp_expiry,omitempty"`
}

type OTPEntry struct {
	Code      string
	Company   string
	ExpiresAt time.Time
}