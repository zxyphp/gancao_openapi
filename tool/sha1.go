package tool

import (
	"crypto/sha1"
	"encoding/hex"
)

// Sha1Signature SHA1签名
func Sha1Signature(data string) string {
	// 创建一个 SHA1 哈希对象
	h := sha1.New()

	// 写入数据
	h.Write([]byte(data))

	// 计算哈希并返回字节数组
	hashBytes := h.Sum(nil)

	// 将字节数组转换为十六进制字符串
	return hex.EncodeToString(hashBytes)
}
