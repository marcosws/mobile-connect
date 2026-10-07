# Mobile Connect

Open-source platform for managing and remotely controlling Android devices through ADB.

Mobile Connect provides an HTTP API and WebSocket endpoints for device management, APK management, remote control, screen streaming, and Appium integration.

## Features

### Device Management

* Discover connected Android devices
* Retrieve device information
* Execute shell commands through ADB

### APK & Application Management

* Install APK files
* Uninstall applications
* List installed packages
* Launch applications
* Stop applications

### Remote Control

* Screen streaming
* Touch input
* Swipe gestures
* Long press
* Text input

### Mobile Automation

* Appium gateway
* W3C WebDriver session creation
* Device selection by Mobile Connect device ID
* Selenium WebDriver compatibility
* Appium Inspector support

## Architecture

Mobile Connect is designed as a lightweight backend between Android devices, automation tools, and user interfaces.

```text
                     ┌──────────────────────┐
                     │   Vue Frontend       │
                     └──────────┬───────────┘
                                │
                         HTTP / WebSocket
                                │
                                ▼
                     ┌──────────────────────┐
                     │   Mobile Connect    │
                     │      Go API          │
                     ├──────────────────────┤
                     │ Device Service       │
                     │ APK Service          │
                     │ Stream Service       │
                     │ Shell Service       │
                     │ Appium Gateway       │
                     └───────┬──────────────┘
                             │
                    ┌────────┴────────┐
                    │                 │
                   ADB              Appium
                    │                 │
                    └────────┬────────┘
                             │
                             ▼
                    Android Devices
```

The backend is implemented in Go using the standard `net/http` package, with ADB providing the communication layer with Android devices.

WebSocket connections are used for screen streaming and remote control.

## Requirements

* Go 1.25+
* Android Debug Bridge (ADB)
* Android device or emulator
* USB debugging enabled and authorized

Appium is optional and only required when using mobile automation features.

## Quick Start

Clone the repository:

```bash
git clone https://github.com/marcosws/mobile-connect.git
cd mobile-connect
```

The repository contains the Go API in [`backend/`](./backend/) and the Vue 3
frontend in [`frontend/`](./frontend/).

Make sure ADB is available:

```bash
adb devices
```

Start Mobile Connect:

```bash
cd backend
go run ./cmd/api
```

The API will be available at:

```text
http://localhost:8080
```

Check connected devices:

```bash
curl http://localhost:8080/devices
```

To start the Vue development server, open another terminal:

```bash
cd frontend
npm install
npm run dev
```

## API

### Devices

| Method | Endpoint              | Description             |
| ------ | --------------------- | ----------------------- |
| GET    | `/devices`            | List connected devices  |
| GET    | `/devices/{id}`       | Get device information  |
| POST   | `/devices/{id}/shell` | Execute a shell command |

### Remote Control

| Method | Endpoint                       | Description                |
| ------ | ------------------------------ | -------------------------- |
| GET    | `/devices/{id}/stream/video`   | Screen streaming WebSocket |
| GET    | `/devices/{id}/stream/control` | Remote control WebSocket   |

### Applications

| Method | Endpoint                  | Description              |
| ------ | ------------------------- | ------------------------ |
| POST   | `/devices/{id}/install`   | Install an APK           |
| POST   | `/devices/{id}/uninstall` | Uninstall an application |
| GET    | `/devices/{id}/packages`  | List installed packages  |
| POST   | `/devices/{id}/launch`    | Launch an application    |
| POST   | `/devices/{id}/stop`      | Stop an application      |

### Appium

| Method                                         | Endpoint                       | Description                |
| ---------------------------------------------- | ------------------------------ | -------------------------- |
| GET                                            | `/appium/status`               | Check Appium availability  |
| POST                                           | `/devices/{id}/appium/session` | Create an Appium session   |
| `/devices/{id}/appium/session/{sessionId}/...` |                                | Execute WebDriver commands |

For complete Appium configuration, client examples, and troubleshooting, see [backend/APPIUM.MD](./backend/APPIUM.MD).

## Appium Integration

Mobile Connect can act as an Appium gateway.

When creating an Appium session through a device-specific endpoint, the Android device is selected using the `{id}` from the Mobile Connect API.

For example:

```text
http://localhost:8080/devices/emulator-5554/appium
```

This allows automation clients to target a specific Android device without relying on device discovery performed by the automation client itself.

The gateway supports Appium clients, Selenium WebDriver and Appium Inspector.

See [backend/APPIUM.MD](./backend/APPIUM.MD) for detailed instructions.

## Configuration

### Appium

Mobile Connect can automatically start Appium when required.

| Variable            | Default                 | Description                |
| ------------------- | ----------------------- | -------------------------- |
| `APPIUM_URL`        | `http://127.0.0.1:4723` | Appium server URL          |
| `APPIUM_COMMAND`    | `appium`                | Appium executable          |
| `APPIUM_AUTO_START` | `true`                  | Automatically start Appium |

Example:

```bash
cd backend
APPIUM_URL=http://127.0.0.1:4723 go run ./cmd/api
```

To use an externally managed Appium server:

```bash
cd backend
APPIUM_AUTO_START=false go run ./cmd/api
```

## Project Structure

```text
mobile-connect/
├── backend/
│   ├── cmd/api/          # Go API entry point
│   ├── internal/         # Backend packages
│   ├── uploads/          # Uploaded APKs
│   ├── storage/          # Local runtime data (ignored by Git)
│   ├── APPIUM.MD         # Appium documentation
│   ├── go.mod
│   └── go.sum
├── frontend/             # Vue 3 + Vite application
├── LICENSE
└── README.MD
```

## Roadmap

### Device & Application Management

* [x] Device discovery
* [x] Device information
* [x] Shell commands
* [x] APK installation
* [x] APK uninstallation
* [x] Application launch
* [x] Application stop
* [ ] APK repository
* [ ] Device groups

### Remote Control

* [x] Screen streaming
* [x] Remote control
* [ ] Recording
* [ ] Logcat

### Automation

* [x] Appium gateway
* [x] Selenium WebDriver compatibility
* [ ] Authentication
* [ ] Multi-user support

### Platform

* [ ] Dashboard improvements
* [ ] Docker support
* [ ] iOS support

## Security

Mobile Connect is currently designed for local or trusted networks.

The current version does not provide authentication or authorization.

Because the API can execute commands and control connected Android devices, do not expose the service directly to the public Internet.

Authentication and multi-user support are planned for future releases.

## Contributing

Contributions, bug reports and suggestions are welcome.

Before opening a pull request, please make sure the project builds successfully and that the proposed change is consistent with the project's architecture and goals.

## License

Mobile Connect is licensed under the Apache License 2.0.

See [LICENSE](./LICENSE) for details.
