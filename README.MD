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
- Criar e controlar sessões de automação Appium para um dispositivo específico.

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
- `GET /appium/status` - verifica a disponibilidade do servidor Appium configurado.
- `POST /devices/{id}/appium/sessions` - cria uma sessão Appium W3C e fixa `appium:udid` no dispositivo `{id}`.
- `POST /devices/{id}/appium/session` e `/devices/{id}/appium/session/{sessionId}/...` - rota WebDriver compatível com Selenium e Appium Inspector.
- `/appium/{path...}` - encaminha comandos WebDriver para o Appium (por exemplo, `POST /appium/session/{sessionId}/element`).

## Como executar

1. Instale o `adb` no sistema.
2. Conecte um dispositivo Android via USB ou rede.
3. Execute:

```bash
go run ./cmd/api
```

4. O servidor ficará disponível em `http://localhost:8080`.

### Appium

O Mobile Connect inicia o Appium automaticamente em `http://127.0.0.1:4723` se
ele ainda não estiver em execução. Instale antes o binário e o driver desejado,
por exemplo `appium driver install uiautomator2`. Para outra porta ou uma base
como `/wd/hub`, defina `APPIUM_URL`:

```bash
APPIUM_URL=http://127.0.0.1:4723/wd/hub go run ./cmd/api
```

Defina `APPIUM_COMMAND` para informar outro caminho do executável. Para conectar
a um Appium remoto ou gerenciado externamente, use `APPIUM_AUTO_START=false`.

Crie a sessão pelo Mobile Connect; o `appium:udid` é sempre substituído pelo ID
do dispositivo informado na URL:

```bash
curl -X POST http://localhost:8080/devices/emulator-5554/appium/sessions \
  -H 'Content-Type: application/json' \
  -d '{"capabilities":{"alwaysMatch":{"platformName":"Android","appium:automationName":"UiAutomator2"}}}'
```

Use o `sessionId` retornado nos demais comandos WebDriver via `/appium`, por
exemplo `DELETE /appium/session/{sessionId}`.

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
