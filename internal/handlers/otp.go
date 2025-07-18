package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"posServer/internal/models"
	"posServer/internal/storage"
	"posServer/internal/utils"
)

func GenerateOTP(w http.ResponseWriter, r *http.Request) {
	company := r.URL.Query().Get("company")
	client := r.URL.Query().Get("client")

	if company == "" || client == "" {
		http.Error(w, "Missing company or client", http.StatusBadRequest)
		return
	}

	// otp := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000) // Simple 6-digit OTP
	expiry := time.Now().Add(5 * time.Minute)

	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		fmt.Println("Error when generating OTP")
		return
	}
	otp := fmt.Sprintf("%06d", n.Int64())

	// Set or update binding with OTP
	storage.Bindings[client] = models.Binding{
		Company:   company,
		ClientIP:  "", // not registered yet
		OTP:       otp,
		OTPExpiry: expiry,
	}

	err = utils.SendOTP(otp, client)
	if err != nil {
		fmt.Printf("Error sending email to %s.\n", client)
		return
	}
	// Normally you'd send this via email/WhatsApp, but for testing:
	fmt.Fprintf(w, "OTP for client '%s' is: %s\n", client, otp)
}
