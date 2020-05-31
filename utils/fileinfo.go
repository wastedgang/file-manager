package utils

import (
	"github.com/mozillazg/go-pinyin"
	"os"
)

func init() {
	pinyinArgs = pinyin.NewArgs()
}

var (
	pinyinArgs pinyin.Args
)

func GetFileSize(filepath string) (int64, error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func GetPinYin(r rune) string {
	result := pinyin.SinglePinyin(r, pinyinArgs)
	if len(result) == 0 {
		return ""
	}
	return result[0]
}
