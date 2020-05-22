package utils

import "os"

// IsExist 文件或文件夹是否存在
func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil || os.IsExist(err)
}

// IsFile 文件是否存在
func IsFile(filepath string) bool {
	stat, err := os.Stat(filepath)
	if err != nil || stat == nil {
		return false
	}
	return !stat.IsDir()
}

// IsDir 文件夹是否存在
func IsDir(filepath string) bool {
	stat, err := os.Stat(filepath)
	if err != nil || stat == nil {
		return false
	}
	return stat.IsDir()
}
