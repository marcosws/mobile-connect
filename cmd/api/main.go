package main

import (
	"log"
	"mobile-connect/internal/adb"
	"mobile-connect/internal/api"
	"mobile-connect/internal/devices"
	"mobile-connect/internal/middleware"
	"mobile-connect/internal/shell"
	"mobile-connect/internal/stream"
	"net/http"
)

func main() {

	adbClient := adb.New()

	deviceService := devices.NewService(adbClient)
	shellService := shell.NewService(adbClient)
	streamService := stream.NewService()

	handler := api.NewHandler(
		deviceService,
		shellService,
		streamService,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /devices", handler.GetDevices)

	mux.HandleFunc("GET /devices/{id}", handler.GetDeviceByID)

	mux.HandleFunc("POST /devices/{id}/shell", handler.Shell)

	mux.HandleFunc(
		"GET /devices/{id}/stream/video",
		handler.StreamVideo,
	)

	mux.HandleFunc(
		"GET /devices/{id}/stream/control",
		handler.StreamControl,
	)

	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

	log.Println("server running on :8080")

	log.Fatal(
		http.ListenAndServe(
			":8080",
			middleware.Cors(mux),
		),
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
