package main

import (
	"fmt"
	"log"
	"net/http"

	"posServer/internal/config"
	"posServer/internal/handlers"
	"posServer/internal/storage"
)

func main() {
	storage.LoadBindings()
	storage.LoadDNSRegistry()

	conf := config.NewConfig()

	http.HandleFunc("/bind", handlers.RegisterClient)
	http.HandleFunc("/route", handlers.RouteHandler)
	http.HandleFunc("/registerCompany", handlers.RegisterCompany)
	http.HandleFunc("/generateOTP", handlers.GenerateOTP)

	lanIP := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)
	// lanIP := "0.0.0.0:3000"                                   // Binds to all interfaces
	fmt.Println("Server started at " + lanIP) // Replace manually or get dynamically
	log.Fatal(http.ListenAndServe(lanIP, nil))
}
