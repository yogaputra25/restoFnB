package handler

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/resto-fnb/backend/internal/middleware"
	"github.com/resto-fnb/backend/pkg/response"
)

func newUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

var allowedExtensions = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/jpg":  ".jpg",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "no file provided")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	ext, ok := allowedExtensions[contentType]
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid file type: only PNG, JPG, GIF, WebP allowed")
		return
	}

	chainDir := filepath.Join("uploads", user.ChainID)
	if err := os.MkdirAll(chainDir, 0755); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create upload directory")
		return
	}

	filename := newUUID() + ext
	dst, err := os.Create(filepath.Join(chainDir, filename))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	url := fmt.Sprintf("/uploads/%s/%s", user.ChainID, filename)
	response.Created(w, map[string]string{"url": url})
}

func UploadRoutes() (string, http.Handler) {
	fs := http.FileServer(http.Dir("uploads"))
	handler := http.StripPrefix("/uploads/", fs)
	return "/uploads/", handler
}
