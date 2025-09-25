package util

import (
	"os"
	"path/filepath"
	"strings"
)

// GetExecutableDir 获取项目根目录下的config路径
func GetExecutableDir() string {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	if strings.Contains(exeDir, "bin") {
		return exeDir
	}

	// 通过工作目录推导（适合go run调试）
	wd, _ := os.Getwd()
	return wd
}
