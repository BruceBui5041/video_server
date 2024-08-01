package utils

import "path/filepath"

func RenameFile(originalName, newBaseName string) string {
	ext := filepath.Ext(originalName)
	return newBaseName + ext
}
