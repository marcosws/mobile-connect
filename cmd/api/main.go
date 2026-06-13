package main

import (
	"log"
	"mobile-connect/internal/adb"
	"mobile-connect/internal/api"
	"mobile-connect/internal/devices"
	"mobile-connect/internal/shell"
	"net/http"
)

func main() {

	adbClient := adb.New()

	deviceService := devices.NewService(adbClient)
	shellService := shell.NewService(adbClient)

	handler := api.NewHandler(
		deviceService,
		shellService,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /devices", handler.GetDevices)

	mux.HandleFunc("GET /devices/{id}", handler.GetDeviceByID)

	mux.HandleFunc("POST /devices/{id}/shell", handler.Shell)

	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

	log.Println("server running on :8080")

	log.Fatal(
		http.ListenAndServe(":8080", mux),
	)

}

/*
	List Devices

	Endpoint: GET /devices

	curl http://localhost:8080/devices

	Response:
	[
		{
			"id": "52000e5003ef25ab",
			"status": "device",
			"manufacturer": "samsung",
			"model": "SM-J701MT",
			"android_version": "8.1.0",
			"sdk_version": "27"
		}
	]

	Get Device by ID

	Endpoint: GET /devices/{id}

	curl http://localhost:8080/devices/52000e5003ef25ab

	Response:
	device: 52000e5003ef25ab


	Shell

	Endpoint: POST /devices/{id}/shell

	curl -X POST http://localhost:8080/devices/52000e5003ef25ab/shell \
	-H "Content-Type: application/json" \
	-d '{
	"command": "getprop ro.product.model"
	}'

	Response:
	{
		"output":"SM-J701MT\n"
	}
*/
