package handlers

import (
	"fmt"
	"net/http"
	"posServer/internal/storage"
	"posServer/internal/utils"
)

func RouteHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := utils.GetClientIP(r)
	fmt.Printf("Incoming route request from IP: %s\n", clientIP)

	company := r.URL.Query().Get("company")
	client := r.URL.Query().Get("client")
	if client == "" || company == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	binding, exists := storage.Bindings[client]
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	// Check if the company matches the binding
	if binding.Company != company {
		http.Error(w, "Company mismatch", http.StatusForbidden)
		return
	}

	// fmt.Printf("Expected IP: %s, Actual IP: %s\n", binding.ClientIP, clientIP)

	// ✅ Compare IP address
	if binding.ClientIP != clientIP {
		http.Error(w, "IP mismatch. Unauthorized access from a different machine.", http.StatusUnauthorized)
		return
	}

	targetURL, ok := storage.DNSRegistry[binding.Company]
	if !ok {
		http.Error(w, "Company not registered", http.StatusNotFound)
		return
	}

	// Redirect client to company database
	http.Redirect(w, r, targetURL, http.StatusFound)
}