# Mobile Connect

Mobile Connect é uma aplicação em Go para gerenciamento e controle de dispositivos Android via ADB.
O projeto fornece uma API HTTP para listar dispositivos, executar comandos shell remotos, instalar e desinstalar APKs,
listar pacotes, iniciar e parar apps, além de capturar a tela e controlar o dispositivo remotamente.

## Visão Geral

- Backend em Go usando `net/http`.
- Integração com ADB para comando remoto em dispositivos Android.
- Websocket para streaming de vídeo e controle remoto.
- Frontend estático servido a partir de `web/index.html`.
- Middleware CORS configurado para `http://localhost:5173`.

## Funcionalidades principais

- Listar dispositivos Android conectados.
- Consultar informações de um dispositivo específico.
- Executar comandos shell remotos via ADB.
- Capturar a tela do dispositivo.
- Enviar toques, swipes, long press e texto.
- Instalar e desinstalar APKs.
- Listar pacotes instalados no dispositivo.
- Iniciar e parar aplicativos.

## Estrutura de pastas

```
mobile-connect/
├── cmd/
│   └── api/
│       └── main.go           # ponto de entrada do servidor HTTP
├── internal/
│   ├── adb/                  # abstração do cliente ADB
│   │   ├── adb.go
│   │   └── client.go
│   ├── api/                  # handlers HTTP e rotas da API
│   │   ├── apk.go
│   │   ├── devices.go
│   │   ├── handler.go
│   │   ├── shell.go
│   │   └── stream.go
│   ├── apk/                  # lógica de instalação/desinstalação e apps
│   │   ├── entity/
│   │   │   ├── apprequest.go
│   │   │   └── uninstall.go
│   │   └── service.go
│   ├── devices/              # serviço de listagem e consulta de dispositivos
│   │   ├── entity/
│   │   │   ├── device.go
│   │   │   └── shell.go
│   │   └── service.go
│   ├── logger/               # logger compartilhado
│   │   └── logger.go
│   ├── middleware/           # middlewares HTTP, por exemplo CORS
│   │   └── cors.go
│   ├── shell/                # serviço de execução de comandos shell via ADB
│   │   └── service.go
│   └── stream/               # captura de tela e controle remoto via ADB
│       ├── service.go
│       └── types.go
├── storage/
│   ├── apks/
│   ├── logs/
│   ├── recordings/
│   ├── screenshots/
│   └── uploads/
├── web/
│   └── index.html            # frontend estático
├── go.mod
├── go.sum
└── README.MD
```

## Endpoints principais

- `GET /devices` - lista dispositivos conectados.
- `GET /devices/{id}` - obtém dados do dispositivo por ID.
- `POST /devices/{id}/shell` - executa um comando shell no dispositivo.
- `GET /devices/{id}/stream/video` - websocket para streaming de tela.
- `GET /devices/{id}/stream/control` - websocket para controle remoto.
- `POST /devices/{id}/install` - instala um APK enviado como multipart/form-data.
- `POST /devices/{id}/uninstall` - desinstala um pacote pelo nome.
- `GET /devices/{id}/packages` - lista pacotes instalados.
- `POST /devices/{id}/launch` - inicia um aplicativo por pacote.
- `POST /devices/{id}/stop` - para um aplicativo em execução.

## Como executar

1. Instale o `adb` no sistema.
2. Conecte um dispositivo Android via USB ou rede.
3. Execute:

```bash
go run ./cmd/api
```

4. O servidor ficará disponível em `http://localhost:8080`.

## Dependências

- `github.com/gorilla/websocket`


Roadmap

[x] Device discovery

[x] Shell

[x] Install APK

[x] Uninstall APK

[x] Launch App

[x] Stop App

[x] Streaming

[x] Remote Control

[ ] APK Repository

[ ] Device Groups

[ ] Recording

[ ] Logcat

[ ] Appium Integration

[ ] Selenium Integration

[ ] Authentication

[ ] Multi-user

[ ] iOS Support

[ ] Dashboard

[ ] Docker

[ ] Kubernetes