package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// generateKey 生成AES加密/解密关键字
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// AESDecrypt AES解密
func AESDecrypt(encrypted, key []byte) (decrypted []byte, resultErr error) {
	defer func() {
		err, ok := recover().(error)
		if err == nil {
			return
		}
		if ok {
			resultErr = err
		} else {
			resultErr = errors.New("decryption error")
		}
		if resultErr != nil {
			decrypted = nil
		}
	}()

	key = generateKey(key)
	block, err := aes.NewCipher(key) // 获取block块
	if err != nil {
		resultErr = err
		return
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decrypted = make([]byte, len(encrypted))
	blockMode.CryptBlocks(decrypted, encrypted)
	decrypted = _PKCS7UnPadding(decrypted)
	return decrypted, nil
}

// AESEncrypt AES加密
func AESEncrypt(originalData, key []byte) (encrypted []byte, resultErr error) {
	defer func() {
		err, ok := recover().(error)
		if err == nil {
			return
		}
		if ok {
			resultErr = err
		} else {
			resultErr = errors.New("encryption error")
		}
		if resultErr != nil {
			encrypted = nil
		}
	}()

	key = generateKey(key)
	block, err := aes.NewCipher(key) // 获取block块
	if err != nil {
		resultErr = err
		return
	}

	originalData = _PKCS7Padding(originalData, block.BlockSize())       // 补码
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()]) //加密模式
	encrypted = make([]byte, len(originalData))                         // 创建明文长度的数组
	blockMode.CryptBlocks(encrypted, originalData)                      // 加密明文
	return encrypted, nil
}

// _PKCS7Padding 补码
func _PKCS7Padding(origData []byte, blockSize int) []byte {
	padding := blockSize - len(origData)%blockSize          // 计算需要补几位数
	padtext := bytes.Repeat([]byte{byte(padding)}, padding) // 在切片后面追加char数量的byte(char)
	return append(origData, padtext...)
}

// _PKCS7UnPadding 去补码
func _PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}
