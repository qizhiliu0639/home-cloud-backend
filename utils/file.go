package utils

import (
	"path/filepath"
	"strings"
)

func GetFileTypeByName(name string) (fileType string) {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".py":
		fallthrough
	case ".go":
		fallthrough
	case ".js":
		fallthrough
	case ".html":
		fallthrough
	case ".css":
		fallthrough
	case ".c":
		fallthrough
	case ".cpp":
		fallthrough
	case ".java":
		fallthrough
	case ".txt":
		fileType = "txt"
	case ".md":
		fileType = "md"
	case ".mp3":
		fallthrough
	case ".wav":
		fallthrough
	case ".flac":
		fallthrough
	case ".wma":
		fallthrough
	case ".ape":
		fallthrough
	case ".aac":
		fileType = "audio"
	case ".mp4":
		fallthrough
	case ".avi":
		fallthrough
	case ".mpg":
		fallthrough
	case ".wmv":
		fallthrough
	case ".mkv":
		fallthrough
	case ".flv":
		fallthrough
	case ".mov":
		fileType = "video"
	case ".jpg":
		fallthrough
	case ".jpeg":
		fallthrough
	case ".png":
		fallthrough
	case ".svg":
		fallthrough
	case ".gif":
		fallthrough
	case ".webp":
		fallthrough
	case ".heic":
		fallthrough
	case ".bmp":
		fileType = "image"
	case ".zip":
		fallthrough
	case ".7z":
		fallthrough
	case ".rar":
		fallthrough
	case ".tar":
		fallthrough
	case ".gz":
		fileType = "zip"
	case ".doc":
		fallthrough
	case ".docx":
		fileType = "doc"
	case ".ppt":
		fallthrough
	case ".pptx":
		fileType = "ppt"
	case ".xls":
		fallthrough
	case ".xlsx":
		fileType = "xls"
	case ".pdf":
		fileType = "pdf"
	case ".exe":
		fileType = "exe"
	default:
		fileType = "other"
	}
	return
}
