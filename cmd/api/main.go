package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mobile-connect/internal/adb"
	"mobile-connect/internal/api"
	"mobile-connect/internal/apk"
	"mobile-connect/internal/appium"
	"mobile-connect/internal/apps"
	"mobile-connect/internal/database"
	"mobile-connect/internal/devices"
	"mobile-connect/internal/middleware"
	"mobile-connect/internal/shell"
	"mobile-connect/internal/stream"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	adbClient := adb.New()
	appiumURL := os.Getenv("APPIUM_URL")
	if appiumURL == "" {
		appiumURL = "http://127.0.0.1:4723"
	}
	appiumClient, err := appium.New(appiumURL)
	if err != nil {
		log.Fatal(err)
	}
	if autoStart, err := appiumAutoStart(); err != nil {
		log.Fatal(err)
	} else if autoStart {
		appiumCommand := os.Getenv("APPIUM_COMMAND")
		if appiumCommand == "" {
			appiumCommand = "appium"
		}
		if err := appium.StartLocal(ctx, appiumClient, appiumCommand); err != nil {
			log.Fatal(err)
		}
	}
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
		appiumClient,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /devices", handler.GetDevices)

	mux.HandleFunc("GET /devices/{id}", handler.GetDeviceByID)

	mux.HandleFunc("POST /devices/{id}/shell", handler.Shell)

	mux.HandleFunc("GET /appium/status", handler.AppiumStatus)

	mux.HandleFunc("POST /devices/{id}/appium/sessions", handler.CreateAppiumSession)

	// This route uses the WebDriver convention expected by Selenium and Appium Inspector.
	mux.HandleFunc("POST /devices/{id}/appium/session", handler.CreateAppiumSession)

	mux.HandleFunc("/devices/{id}/appium/{path...}", handler.AppiumDeviceProxy)

	mux.HandleFunc("/appium/{path...}", handler.AppiumProxy)

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

	server := &http.Server{Addr: ":8080", Handler: middleware.Cors(mux)}
	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}

func appiumAutoStart() (bool, error) {
	value := os.Getenv("APPIUM_AUTO_START")
	if value == "" {
		return true, nil
	}
	enabled, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid APPIUM_AUTO_START value %q: use true or false", value)
	}
	return enabled, nil
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
