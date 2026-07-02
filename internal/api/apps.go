package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"mobile-connect/internal/apps/entity"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var (
	rePackage   = regexp.MustCompile(`package: name='([^']+)' versionCode='(\d+)' versionName='([^']*)'`)
	reMinSDK    = regexp.MustCompile(`sdkVersion:'(\d+)'`)
	reTargetSDK = regexp.MustCompile(`targetSdkVersion:'(\d+)'`)
)

func (h *Handler) CreateApp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		http.Error(w, "erro ao parsear multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("apk")
	if err != nil {
		http.Error(w, "erro ao ler arquivo apk: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	id := uuid.NewString()
	fileName := id + ".apk"
	uploadPath := filepath.Join("./uploads", fileName)
	storagePath := filepath.Join("./storage", "apks", fileName)

	if err := os.MkdirAll(filepath.Dir(storagePath), 0o755); err != nil {
		http.Error(w, "erro ao criar diretório de armazenamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(uploadPath)
	if err != nil {
		http.Error(w, "erro ao salvar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hasher := sha256.New()
	size, err := io.Copy(io.MultiWriter(dst, hasher), file)
	if err != nil {
		dst.Close()
		http.Error(w, "erro ao copiar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := dst.Close(); err != nil {
		http.Error(w, "erro ao fechar arquivo de upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	uploadFile, err := os.Open(uploadPath)
	if err != nil {
		http.Error(w, "erro ao abrir arquivo de upload: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadFile.Close()

	if _, err := uploadFile.Seek(0, 0); err != nil {
		http.Error(w, "erro ao reiniciar arquivo temporário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	storageFile, err := os.Create(storagePath)
	if err != nil {
		http.Error(w, "erro ao salvar apk no armazenamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := uploadFile.Seek(0, 0); err != nil {
		storageFile.Close()
		http.Error(w, "erro ao reiniciar arquivo de upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(storageFile, uploadFile); err != nil {
		storageFile.Close()
		http.Error(w, "erro ao copiar apk para armazenamento: "+err.Error(), http.StatusInternalServerError)
		return
	}
	storageFile.Close()

	meta, err := extractApkMetadata(uploadPath)
	if err != nil {
		meta = &apkMetadata{}
	}

	app := entity.App{
		ID:           id,
		FileName:     fileName,
		OriginalName: header.Filename,
		PackageName:  meta.PackageName,
		VersionName:  meta.VersionName,
		VersionCode:  meta.VersionCode,
		MinSDK:       meta.MinSDK,
		TargetSDK:    meta.TargetSDK,
		Size:         size,
		SHA256:       hex.EncodeToString(hasher.Sum(nil)),
		CreatedAt:    time.Now(),
	}

	if err := h.appService.Create(app); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type apkMetadata struct {
	PackageName string
	VersionName string
	VersionCode string
	MinSDK      int
	TargetSDK   int
}

func extractApkMetadata(apkPath string) (*apkMetadata, error) {
	out, err := exec.Command("aapt", "dump", "badging", apkPath).CombinedOutput()
	if err != nil {
		return &apkMetadata{}, nil
	}

	output := string(out)
	meta := &apkMetadata{}

	if m := rePackage.FindStringSubmatch(output); m != nil {
		meta.PackageName = m[1]
		meta.VersionCode = m[2]
		meta.VersionName = m[3]
	}

	if m := reMinSDK.FindStringSubmatch(output); m != nil {
		meta.MinSDK, _ = strconv.Atoi(m[1])
	}

	if m := reTargetSDK.FindStringSubmatch(output); m != nil {
		meta.TargetSDK, _ = strconv.Atoi(m[1])
	}

	return meta, nil
}

func (h *Handler) GetApps(w http.ResponseWriter, r *http.Request) {
	appsList, err := h.appService.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appsList)
}

func (h *Handler) GetAppByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	app, err := h.appService.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

func (h *Handler) DeleteApp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.appService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/** ===========================================*/

func (h *Handler) InstallStoredApp(
	w http.ResponseWriter,
	r *http.Request,
) {

	appID := r.PathValue("appId")
	deviceID := r.PathValue("deviceId")

	output, err := h.appService.Install(
		appID,
		deviceID,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"device": deviceID,
			"appId":  appID,
			"status": "installed",
			"output": output,
		},
	)
}
