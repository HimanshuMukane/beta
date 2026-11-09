package utils

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ServeSPA serves static files from the given root directory and falls back to index.html
// for client-side routing (React CSR).
func ServeSPA(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean and join requested path
		requestedPath := filepath.Clean(r.URL.Path)
		fullPath := filepath.Join(root, requestedPath)

		// Check if the requested file exists and is not a directory
		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() {
			// Cache hashed assets immutably, revalidate others
			if strings.Contains(requestedPath, "/assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			http.ServeFile(w, r, fullPath)
			return
		}

		// Fallback to index.html for CSR routes (never cache index.html long-term)
		indexPath := filepath.Join(root, "index.html")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.ServeFile(w, r, indexPath)
	})
}

