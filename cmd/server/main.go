package main

import (
	"fmt"
	"log"
	"net/http"

	"posServer/internal/handlers"
	"posServer/internal/storage"
)

func main() {
	storage.LoadBindings()
	storage.LoadDNSRegistry()

	http.HandleFunc("/bind", handlers.RegisterClient)
	http.HandleFunc("/route", handlers.RouteHandler)
	http.HandleFunc("/registerCompany", handlers.RegisterCompany)
	http.HandleFunc("/generateOTP", handlers.GenerateOTP)

	lanIP := "0.0.0.0:3000"                                   // Binds to all interfaces
	fmt.Println("Server started at http://192.168.0.66:3000") // Replace manually or get dynamically
	log.Fatal(http.ListenAndServe(lanIP, nil))
}
