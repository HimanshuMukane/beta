package utils

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ServeSPA serves static files from the given root directory and falls back to index.html
// for client-side routing (React CSR). Automatically serves pre-compressed .br and .gz files.
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

			// Pre-compressed asset negotiation (Brotli > Gzip)
			acceptEnc := r.Header.Get("Accept-Encoding")
			contentType := mimeTypeByExt(filepath.Ext(fullPath))

			if strings.Contains(acceptEnc, "br") {
				if _, errBr := os.Stat(fullPath + ".br"); errBr == nil {
					w.Header().Set("Content-Encoding", "br")
					if contentType != "" {
						w.Header().Set("Content-Type", contentType)
					}
					http.ServeFile(w, r, fullPath+".br")
					return
				}
			}

			if strings.Contains(acceptEnc, "gzip") {
				if _, errGz := os.Stat(fullPath + ".gz"); errGz == nil {
					w.Header().Set("Content-Encoding", "gzip")
					if contentType != "" {
						w.Header().Set("Content-Type", contentType)
					}
					http.ServeFile(w, r, fullPath+".gz")
					return
				}
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

func mimeTypeByExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".js", ".mjs":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".html":
		return "text/html; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".json":
		return "application/json"
	case ".woff2":
		return "font/woff2"
	case ".woff":
		return "font/woff"
	case ".ttf":
		return "font/ttf"
	default:
		return ""
	}
}

