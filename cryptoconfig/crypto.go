package cryptoconfig

import (
	"encoding/json"
	"github.com/farseer810/file-manager/utils"
)

// DecryptConfig 解密配置信息
func DecryptConfig(encryptedData, key []byte) (*Config, error) {
	decryptedBytes, err := utils.AESDecrypt(encryptedData, key)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = json.Unmarshal(decryptedBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// EncryptConfig 加密配置信息
func EncryptConfig(config *Config, key []byte) ([]byte, error) {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return utils.AESEncrypt(jsonBytes, key)
}
