package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type GcOpenApi struct {
	Url       string
	AccessKey string
	SecretKey string
}

// NewGcOpenApi 初始化构造函数，接收 GcOpenApi struct 和 isProd 作为参数
func NewGcOpenApi(conf GcOpenApi, isProd bool) *GcOpenApi {
	baseURL := "http://dev-gapis-base.igancao.com/oapi"
	if isProd {
		baseURL = "https://gapis-base-outer.igancao.com/oapi"
	}

	conf.Url = baseURL
	return &conf
}

// ExecApi 执行API请求
func (api *GcOpenApi) ExecApi(pkg, class string, inParam map[string]interface{}) (map[string]interface{}, error) {
	inParam["package"] = pkg
	inParam["class"] = class
	return api.transmit(inParam)
}

// 数据传输层
func (api *GcOpenApi) transmit(data map[string]interface{}) (map[string]interface{}, error) {
	timestamp := time.Now().Unix()
	noise := api.randomString(8)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %v", err)
	}

	signature := api.sha1Signature(fmt.Sprintf("%s%d%s%s", jsonData, timestamp, noise, api.SecretKey))
	encryptedData, err := api.encrypt(string(jsonData), api.SecretKey)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", api.Url, strings.NewReader(encryptedData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("AK", api.AccessKey)
	req.Header.Set("Signature", signature)
	req.Header.Set("UTC-Timestamp", fmt.Sprint(timestamp))
	req.Header.Set("NOISE", noise)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解密响应
	decryptedData, err := api.decrypt(string(body), api.SecretKey)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(decryptedData), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SHA1签名
func (api *GcOpenApi) sha1Signature(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// AES加密
func (api *GcOpenApi) encrypt(data, key string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(data))
	ecb := newECBEncrypter(block)
	ecb.CryptBlocks(ciphertext, []byte(data))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AES解密
func (api *GcOpenApi) decrypt(data, key string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}

	encryptedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %v", err)
	}

	ciphertext := make([]byte, len(encryptedData))
	ecb := newECBDecrypter(block)
	ecb.CryptBlocks(ciphertext, encryptedData)

	return string(ciphertext), nil
}

// 新建 AES Cipher Block
func newCipherBlock(key string) (cipher.Block, error) {
	return aes.NewCipher([]byte(key[:16]))
}

// ECB加密器
type ecbEncrypter struct {
	b cipher.Block
}

func newECBEncrypter(b cipher.Block) cipher.BlockMode {
	return &ecbEncrypter{b}
}

func (x *ecbEncrypter) BlockSize() int { return x.b.BlockSize() }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.BlockSize() != 0 {
		panic("crypto/aes: input not full blocks")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.BlockSize()])
		src = src[x.BlockSize():]
		dst = dst[x.BlockSize():]
	}
}

// ECB解密器
type ecbDecrypter struct {
	b cipher.Block
}

func newECBDecrypter(b cipher.Block) cipher.BlockMode {
	return &ecbDecrypter{b}
}

func (x *ecbDecrypter) BlockSize() int { return x.b.BlockSize() }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.BlockSize() != 0 {
		panic("crypto/aes: input not full blocks")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.BlockSize()])
		src = src[x.BlockSize():]
		dst = dst[x.BlockSize():]
	}
}

// 获取随机字符串
func (api *GcOpenApi) randomString(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
