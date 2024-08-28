package tool

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
)

func pKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Encrypt AES加密（ECB模式）
func Encrypt(plainText, key string) (string, error) {
	// Ensure the key is 16 bytes for AES-128
	if len(key) > 16 {
		key = key[:16]
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plainTextBytes := pKCS7Padding([]byte(plainText), block.BlockSize())

	cipherText := make([]byte, len(plainTextBytes))
	mode := NewECBEncrypter(block) // Ensure ECB mode
	mode.CryptBlocks(cipherText, plainTextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt AES解密（ECB模式）
func Decrypt(data, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key[:16]))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	plainText := make([]byte, len(cipherText))
	mode := NewECBDecrypter(block)
	mode.CryptBlocks(plainText, cipherText)

	// 去除PKCS7填充
	padding := int(plainText[len(plainText)-1])
	return string(plainText[:len(plainText)-padding]), nil
}
