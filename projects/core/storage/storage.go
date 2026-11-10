package storage

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type StorageProvider interface {
	GenerateUploadURL(ctx context.Context, key string, contentType string, expires time.Duration) (string, error)
	GenerateDownloadURL(ctx context.Context, key string, fileName string, expires time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
	ObjectExists(ctx context.Context, key string) (bool, error)
}

type LocalStorageProvider struct {
	BaseDir   string
	SignKey   []byte
	PublicURL string
}

func NewLocalStorageProvider(baseDir string, secret string, publicURL string) *LocalStorageProvider {
	if baseDir == "" {
		baseDir = "./uploads"
	}
	if publicURL == "" {
		publicURL = "/api/storage/files"
	}
	return &LocalStorageProvider{
		BaseDir:   baseDir,
		SignKey:   []byte(secret),
		PublicURL: publicURL,
	}
}

func (s *LocalStorageProvider) generateSignature(key string, exp int64) string {
	mac := hmac.New(sha256.New, s.SignKey)
	mac.Write([]byte(fmt.Sprintf("%s:%d", key, exp)))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *LocalStorageProvider) VerifySignature(key string, expStr string, sig string) bool {
	exp, err := time.ParseDuration(expStr)
	if err != nil {
		var expUnix int64
		if _, scanErr := fmt.Sscanf(expStr, "%d", &expUnix); scanErr != nil || time.Now().Unix() > expUnix {
			return false
		}
		expected := s.generateSignature(key, expUnix)
		return hmac.Equal([]byte(sig), []byte(expected))
	}
	_ = exp
	return false
}

func (s *LocalStorageProvider) GenerateUploadURL(ctx context.Context, key string, contentType string, expires time.Duration) (string, error) {
	exp := time.Now().Add(expires).Unix()
	sig := s.generateSignature(key, exp)
	return fmt.Sprintf("%s/upload?key=%s&exp=%d&sig=%s&content_type=%s", s.PublicURL, url.QueryEscape(key), exp, sig, url.QueryEscape(contentType)), nil
}

func (s *LocalStorageProvider) GenerateDownloadURL(ctx context.Context, key string, fileName string, expires time.Duration) (string, error) {
	exp := time.Now().Add(expires).Unix()
	sig := s.generateSignature(key, exp)
	return fmt.Sprintf("%s/download?key=%s&filename=%s&exp=%d&sig=%s", s.PublicURL, url.QueryEscape(key), url.QueryEscape(fileName), exp, sig), nil
}

func (s *LocalStorageProvider) DeleteObject(ctx context.Context, key string) error {
	cleanKey := filepath.Clean(key)
	if strings.Contains(cleanKey, "..") {
		return fmt.Errorf("invalid storage key path")
	}
	targetPath := filepath.Join(s.BaseDir, cleanKey)
	return os.Remove(targetPath)
}

func (s *LocalStorageProvider) ObjectExists(ctx context.Context, key string) (bool, error) {
	cleanKey := filepath.Clean(key)
	if strings.Contains(cleanKey, "..") {
		return false, fmt.Errorf("invalid storage key path")
	}
	targetPath := filepath.Join(s.BaseDir, cleanKey)
	info, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil && !info.IsDir(), err
}
