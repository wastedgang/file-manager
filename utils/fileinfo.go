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

func GetFirstPinYinLetter(s string) (r rune) {
	result := pinyin.Pinyin(s, pinyinArgs)
	if len(result) == 0 {
		return
	}
	if len(result[0]) == 0 {
		return
	}
	if len(result[0][0]) == 0 {
		return
	}
	return []rune(result[0][0])[0]
}
