package main

import (
  "fmt"
  "net/http"
  "net/http/httptest"
  "strings"
)

func main() {
  body := strings.NewReader("--abc\r\nContent-Disposition: form-data; name=\"apk\"; filename=\"a.apk\"\r\nContent-Type: application/octet-stream\r\n\r\n123\r\n--abc--\r\n")
  req := httptest.NewRequest(http.MethodPost, "/apps", body)
  req.Header.Set("Content-Type", "multipart/form-data; boundary=abc")
  err := req.ParseMultipartForm(64 << 20)
  fmt.Printf("parse err=%v\n", err)
  if err == nil {
    f, _, err2 := req.FormFile("apk")
    fmt.Printf("formfile err=%v\n", err2)
    if f != nil { _ = f.Close() }
  }
}
