package client

import (
	"encoding/json"
	"fmt"
	"github.com/zxyphp/gancao_openapi/tool"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// GcOpenApi 结构体定义
type GcOpenApi struct {
	Url       string
	AccessKey string
	SecretKey string
}

// SortedData 定义一个结构体来保持顺序
type SortedData struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Package string `json:"package"`
	Class   string `json:"class"`
}

// NewGcOpenApi 初始化构造函数
func NewGcOpenApi(conf GcOpenApi, isProd bool) *GcOpenApi {
	baseURL := "http://dev-gapis-base.igancao.com/oapi"
	if isProd {
		baseURL = "https://gapis-base-outer.igancao.com/oapi"
	}
	conf.Url = baseURL
	return &conf
}

// ApiRequest 用于封装 API 请求的参数
type ApiRequest struct {
	Package string
	Class   string
	Params  map[string]interface{}
}

// ExecApi 执行API请求
func (api *GcOpenApi) ExecApi(req ApiRequest) (map[string]interface{}, error) {
	req.Params["package"] = req.Package
	req.Params["class"] = req.Class

	return api.transmit(req.Params)
}

// 数据传输层
func (api *GcOpenApi) transmit(data map[string]interface{}) (map[string]interface{}, error) {
	timestamp := time.Now().Unix()
	noise := tool.RandStr(8)

	// 创建SortedData实例并填充数据
	sortedData := SortedData{
		Name:    data["name"].(string),
		Age:     data["age"].(int),
		Package: data["package"].(string),
		Class:   data["class"].(string),
	}

	// 将结构体转换为JSON格式的字符串
	jsonData, err := json.Marshal(sortedData)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %v", err)
	}
	signature := tool.Sha1Signature(string(jsonData) + fmt.Sprint(timestamp) + noise + api.SecretKey)
	encryptedData, err := tool.Encrypt(string(jsonData), api.SecretKey)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", api.Url, strings.NewReader(encryptedData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Connection", "close")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("AK", api.AccessKey)
	req.Header.Set("Signature", signature)
	req.Header.Set("UTC-Timestamp", fmt.Sprint(timestamp))
	req.Header.Set("NOISE", noise)
	req.Header.Set("Expect", "")

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
	decryptedData, err := tool.Decrypt(string(body), api.SecretKey)
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
