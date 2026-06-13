

Mobile Connect

mobile-connect/
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── api/              # handlers HTTP / websocket
│   ├── devices/          # regras de device + adb wrapper
│   ├── stream/           # abstração de streaming (scrcpy/ffmpeg/adb)
│   ├── scrcpy/           # integração scrcpy (process manager)
│   ├── adb/              # wrapper adb puro
│   ├── ffmpeg/           # wrapper ffmpeg (opcional)
│   ├── logcat/
│   └── recordings/
│
├── pkg/                  # (opcional) libs reutilizáveis
│
├── web/
│   ├── video/
│   └── static/
│
├── scripts/
│   ├── install.sh
│   ├── deps.sh
│
├── Makefile
└── go.mod