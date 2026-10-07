package logger

import (
	"log"
	"os"
)

var Logger = log.New(
	os.Stdout,
	"[MOBILE-CONNECT] ",
	log.Ldate|log.Ltime|log.Lshortfile,
)
