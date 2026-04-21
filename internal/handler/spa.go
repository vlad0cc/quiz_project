package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func NewSPAHandler(distDir string) http.Handler {
	fileServer := http.FileServer(http.Dir(distDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/")
		if path == "." || path == "" {
			http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
			return
		}

		fullPath := filepath.Join(distDir, path)
		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
	})
}
