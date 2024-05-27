package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDirExists 检查目录是否存在，如果不存在则创建它
func EnsureDirExists(dir string) error {
	// 获取当前可执行文件的路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// 获取可执行文件的目录
	execDir := filepath.Dir(execPath)

	// 构建目标目录的绝对路径
	absDir := filepath.Join(execDir, dir)

	// 使用 os.Stat 获取文件信息
	_, err = os.Stat(absDir)
	if os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.MkdirAll(absDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		fmt.Printf("Directory %s created.\n", absDir)
	} else if err != nil {
		// 其他错误
		return fmt.Errorf("failed to check directory: %v", err)
	} else {
		// 目录存在
		fmt.Printf("Directory %s already exists.\n", absDir)
	}
	return nil
}

func main() {
	// 定义相对目录路径
	dir := "./etc/conf"

	// 调用 EnsureDirExists 函数检查并创建目录
	err := EnsureDirExists(dir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
