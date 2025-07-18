package utils

import (
	"net"
	"net/http"
)

func GetClientIP(r *http.Request) string {
	// If behind proxy or load balancer (e.g., nginx), check X-Forwarded-For
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// May contain multiple IPs: client, proxy1, proxy2...
		return forwarded
	}
	// Else, use direct connection IP
	ip := r.RemoteAddr
	// Remove port if exists (e.g., "192.168.1.10:51532")
	if host, _, err := net.SplitHostPort(ip); err == nil {
		return host
	}
	return ip
}