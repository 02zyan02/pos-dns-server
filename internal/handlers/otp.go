package handlers

import (
	"fmt"
	"net/http"
	"time"

	"posServer/internal/models"
	"posServer/internal/storage"
) 

func GenerateOTP(w http.ResponseWriter, r *http.Request) {
	company := r.URL.Query().Get("company")
	client := r.URL.Query().Get("client")

	if company == "" || client == "" {
		http.Error(w, "Missing company or client", http.StatusBadRequest)
		return
	}

	otp := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000) // Simple 6-digit OTP
	expiry := time.Now().Add(5 * time.Minute)

	// Set or update binding with OTP
	storage.Bindings[client] = models.Binding{
		Company:   company,
		ClientIP:  "", // not registered yet
		OTP:       otp,
		OTPExpiry: expiry,
	}

	// Normally you'd send this via email/WhatsApp, but for testing:
	fmt.Fprintf(w, "OTP for client '%s' is: %s\n", client, otp)
}