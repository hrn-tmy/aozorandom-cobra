package cache

import (
	"os"
	"path/filepath"
	"time"
)

const cacheExpiry = 7 * 24 * time.Hour

// CachePath は、キャッシュパスを生成します
func CachePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(dir, "aozora")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	date := time.Now().Format("2006-01-02")

	return filepath.Join(cacheDir, date + "-list.csv"), nil
}

func IsCacheValid(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return time.Since(info.ModTime()) < cacheExpiry
}

func LoadCache(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func SaveData(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}