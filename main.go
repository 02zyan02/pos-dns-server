package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

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

var otpStore = map[string]OTPEntry{}
var otpMutex = sync.Mutex{}

var passkeyBindings = map[string]Binding{}

var dnsRegistry = map[string]string{}

const dnsRegistryFile = "dnsRegistry.json"

const bindingsFile = "bindings.json"

func main() {
	loadBindings()
	loadDNSRegistry()

	http.HandleFunc("/bind", registerClientHandler)
	http.HandleFunc("/route", routeHandler)
	http.HandleFunc("/registerCompany", registerCompanyHandler)
	http.HandleFunc("/generateOTP", generateOTPHandler)

	lanIP := "0.0.0.0:3000"                                   // Binds to all interfaces
	fmt.Println("Server started at http://192.168.0.66:3000") // Replace manually or get dynamically
	log.Fatal(http.ListenAndServe(lanIP, nil))
}

func loadBindings() {
	data, err := ioutil.ReadFile(bindingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return // No file = first run
		}
		log.Fatalf("Failed to read bindings file: %v", err)
	}

	err = json.Unmarshal(data, &passkeyBindings)
	if err != nil {
		log.Fatalf("Failed to parse bindings JSON: %v", err)
	}
}

func loadDNSRegistry() {
	data, err := ioutil.ReadFile(dnsRegistryFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No dnsRegistry.json found, using empty registry.")
			return
		}
		log.Fatalf("Failed to read DNS registry file: %v", err)
	}

	err = json.Unmarshal(data, &dnsRegistry)
	if err != nil {
		log.Fatalf("Failed to parse DNS registry: %v", err)
	}
}

func saveBindings() {
	data, err := json.MarshalIndent(passkeyBindings, "", "  ")
	if err != nil {
		log.Printf("Failed to serialize bindings: %v", err)
		return
	}

	err = ioutil.WriteFile(bindingsFile, data, 0644)
	if err != nil {
		log.Printf("Failed to write bindings file: %v", err)
	}
}

func registerClientHandler(w http.ResponseWriter, r *http.Request) {
	client := r.URL.Query().Get("client")
	company := r.URL.Query().Get("company")
	otp := r.URL.Query().Get("otp")

	if client == "" || company == "" || otp == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	binding, exists := passkeyBindings[client]
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	if binding.Company != company {
		http.Error(w, "Company mismatch", http.StatusForbidden)
		return
	}

	fmt.Printf("binding OTP: '%v'\n", binding)

	if binding.OTP != otp || time.Now().After(binding.OTPExpiry) {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// Update binding with IP, and clear OTP
	binding.ClientIP = getClientIP(r)
	binding.OTP = ""
	binding.OTPExpiry = time.Time{} // zero value
	passkeyBindings[client] = binding

	saveBindings()

	fmt.Fprintf(w, "Client '%s' registered successfully with IP %s\n", client, binding.ClientIP)
}

func registerCompanyHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	company := r.URL.Query().Get("company")
	if company == "" || ip == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	if _, exists := dnsRegistry[company]; exists {
		http.Error(w, "Company already registered", http.StatusConflict)
		return
	}

	// Register new company with a dummy URL
	dnsRegistry[company] = fmt.Sprintf("http://%s", ip)
	saveDNSRegistry()

	fmt.Fprintf(w, "Company '%s' registered successfully.\n", company)

}

func saveDNSRegistry() {
	data, err := json.MarshalIndent(dnsRegistry, "", "  ")
	if err != nil {
		log.Printf("Failed to serialize DNS registry: %v", err)
		return
	}

	err = ioutil.WriteFile(dnsRegistryFile, data, 0644)
	if err != nil {
		log.Printf("Failed to write DNS registry file: %v", err)
	}
}

func getClientIP(r *http.Request) string {
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

func routeHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	fmt.Printf("Incoming route request from IP: %s\n", clientIP)

	company := r.URL.Query().Get("company")
	client := r.URL.Query().Get("client")
	if client == "" || company == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	binding, exists := passkeyBindings[client]
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

	targetURL, ok := dnsRegistry[binding.Company]
	if !ok {
		http.Error(w, "Company not registered", http.StatusNotFound)
		return
	}

	// Redirect client to company database
	http.Redirect(w, r, targetURL, http.StatusFound)
}

func generateOTPHandler(w http.ResponseWriter, r *http.Request) {
	company := r.URL.Query().Get("company")
	client := r.URL.Query().Get("client")

	if company == "" || client == "" {
		http.Error(w, "Missing company or client", http.StatusBadRequest)
		return
	}

	otp := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000) // Simple 6-digit OTP
	expiry := time.Now().Add(5 * time.Minute)

	// Set or update binding with OTP
	passkeyBindings[client] = Binding{
		Company:   company,
		ClientIP:  "", // not registered yet
		OTP:       otp,
		OTPExpiry: expiry,
	}

	// Normally you'd send this via email/WhatsApp, but for testing:
	fmt.Fprintf(w, "OTP for client '%s' is: %s\n", client, otp)
}
