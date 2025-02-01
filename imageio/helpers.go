package imageio

import (
	"path/filepath"
	"strings"
)

func GetFileExtension(filePath string) string {
	return filepath.Ext(filePath)
}

func GetFileName(filePath string) string {
	return filepath.Base(filePath)
}

func GetPureFileName(filePath string) string {
	return strings.TrimSuffix(GetFileName(filePath), GetFileExtension(filePath))
}
