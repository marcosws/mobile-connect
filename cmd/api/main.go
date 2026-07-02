package main

import (
	"log"
	"mobile-connect/internal/adb"
	"mobile-connect/internal/api"
	"mobile-connect/internal/apk"
	"mobile-connect/internal/apps"
	"mobile-connect/internal/database"
	"mobile-connect/internal/devices"
	"mobile-connect/internal/middleware"
	"mobile-connect/internal/shell"
	"mobile-connect/internal/stream"
	"net/http"
)

func main() {

	adbClient := adb.New()
	db := database.New()

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	appRepository := apps.NewRepository(db)

	deviceService := devices.NewService(adbClient)
	shellService := shell.NewService(adbClient)
	apkService := apk.NewService(adbClient)
	streamService := stream.NewService()

	appService := apps.NewService(
		adbClient,
		appRepository,
		apkService,
	)

	handler := api.NewHandler(
		deviceService,
		shellService,
		streamService,
		apkService,
		appService,
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

	mux.HandleFunc(
		"POST /devices/{id}/install",
		handler.InstallAPK,
	)

	mux.HandleFunc(
		"POST /devices/{id}/uninstall",
		handler.UninstallAPK,
	)

	mux.HandleFunc(
		"GET /devices/{id}/packages",
		handler.GetPackages,
	)

	mux.HandleFunc(
		"POST /devices/{id}/launch",
		handler.LaunchApp,
	)

	mux.HandleFunc(
		"POST /devices/{id}/stop",
		handler.StopApp,
	)

	mux.HandleFunc("POST /apps", handler.CreateApp)

	mux.HandleFunc("GET /apps", handler.GetApps)

	mux.HandleFunc("GET /apps/{id}", handler.GetAppByID)

	mux.HandleFunc("DELETE /apps/{id}", handler.DeleteApp)

	mux.HandleFunc(
		"POST /apps/{appId}/devices/{deviceId}/install",
		handler.InstallStoredApp,
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
