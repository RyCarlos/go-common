package fileutil

import (
	"fmt"
	"mime/multipart"
)

// 文件大小单位
const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

// FormatFileSize 格式化文件大小显示
func FormatFileSize(size int64) string {
	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// FormatFileSizePrecise 更精确的版本
func FormatFileSizePrecise(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	var unit string
	value := float64(size)

	for i := 0; i < len(units); i++ {
		unit = units[i]
		if value < 1024.0 || i == len(units)-1 {
			break
		}
		value /= 1024.0
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}

// CheckFileSizeLimit 检查文件大小限制
func CheckFileSizeLimit(header *multipart.FileHeader, maxSize int64) error {
	if header.Size > maxSize {
		return fmt.Errorf("文件大小 %s 超过限制 %s",
			FormatFileSize(header.Size),
			FormatFileSize(maxSize))
	}
	return nil
}
