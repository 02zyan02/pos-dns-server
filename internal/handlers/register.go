package handlers

import(
	"fmt"
	"net/http"
	"posServer/internal/storage"
)

func RegisterCompany(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	company := r.URL.Query().Get("company")
	if company == "" || ip == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	if _, exists := storage.DNSRegistry[company]; exists {
		http.Error(w, "Company already registered", http.StatusConflict)
		return
	}

	storage.DNSRegistry[company] = fmt.Sprintf("http://%s", ip)
	storage.SaveDNSRegistry()

	fmt.Fprintf(w, "Company '%s' registered successfully.\n", company)

}