package main

import (
	"log"
	"mobile-connect/internal/adb"
	"mobile-connect/internal/api"
	"mobile-connect/internal/devices"
	"net/http"
)

func main() {

	adbClient := adb.New()

	deviceService := devices.NewService(adbClient)

	handler := api.NewHandler(
		deviceService,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /devices", handler.GetDevices)

	mux.HandleFunc("GET /devices/{id}", handler.GetDeviceByID)

	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

	log.Println("server running on :8080")

	log.Fatal(
		http.ListenAndServe(":8080", mux),
	)

}
