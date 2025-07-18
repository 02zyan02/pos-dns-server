package handlers

import (
	"fmt"
	"net/http"
	"time"

	"posServer/internal/storage"
	"posServer/internal/utils"
)

func RegisterClient(w http.ResponseWriter, r *http.Request) {
	client := r.URL.Query().Get("client")
	company := r.URL.Query().Get("company")
	otp := r.URL.Query().Get("otp")

	if client == "" || company == "" || otp == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	binding, ok := storage.Bindings[client]
	if !ok || binding.Company != company {
		http.Error(w, "Invalid client or company", http.StatusForbidden)
		return
	} 

	if binding.OTP != otp || time.Now().After(binding.OTPExpiry) {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// Update binding with IP, and clear OTP
	binding.ClientIP = utils.GetClientIP(r)
	binding.OTP = ""
	binding.OTPExpiry = time.Time{} // zero value
	storage.Bindings[client] = binding

	storage.SaveBindings()

	fmt.Fprintf(w, "Client '%s' registered successfully with IP %s\n", client, binding.ClientIP)
}
