package tool

import (
	"crypto/cipher"
)

// ECB 模式的加密器
type ecbEncryptor struct {
	b         cipher.Block
	blockSize int
}

func newECBEncryptor(b cipher.Block) *ecbEncryptor {
	return &ecbEncryptor{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

func (x *ecbEncryptor) BlockSize() int { return x.blockSize }

func (x *ecbEncryptor) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}

	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// ECB 模式的解密器
type ecbDecryptor struct {
	b         cipher.Block
	blockSize int
}

func newECBDecryptor(b cipher.Block) *ecbDecryptor {
	return &ecbDecryptor{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

func (x *ecbDecryptor) BlockSize() int { return x.blockSize }

func (x *ecbDecryptor) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}

	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// NewECBEncrypter 创建 ECB 加密模式
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return newECBEncryptor(b)
}

// NewECBDecrypter 创建 ECB 解密模式
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return newECBDecryptor(b)
}
