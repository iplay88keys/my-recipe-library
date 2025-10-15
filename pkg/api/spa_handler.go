package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type spaHandler struct {
	staticPath string
	indexPath  string
	root       *os.Root
}

func NewSPAHandler(staticPath, indexPath string) (*spaHandler, error) {
	root, err := os.OpenRoot(staticPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open root directory %s: %w", staticPath, err)
	}

	return &spaHandler{
		staticPath: staticPath,
		indexPath:  indexPath,
		root:       root,
	}, nil
}

func (h *spaHandler) Close() error {
	if h.root != nil {
		return h.root.Close()
	}
	return nil
}

// Serves static files or falls back to index.html for SPA routing
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean path and remove leading slash
	cleanPath := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/")

	// Try to serve the requested file
	fileInfo, err := h.root.Stat(cleanPath)
	if err == nil {
		file, err := h.root.Open(cleanPath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)

			return
		}
		defer file.Close()

		http.ServeContent(w, r, cleanPath, fileInfo.ModTime(), file)

		return
	}

	file, err := h.root.Open(h.indexPath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}
	defer file.Close()

	info, err := h.root.Stat(h.indexPath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	// Fallback to serving the index.html for SPA routing
	http.ServeContent(w, r, h.indexPath, info.ModTime(), file)
}
