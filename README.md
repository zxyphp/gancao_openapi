[![Build Status](https://travis-ci.org/nacos-group/nacos-sdk-go.svg?branch=master)](https://travis-ci.org/nacos-group/nacos-sdk-go) [![Go Report Card](https://goreportcard.com/badge/github.com/nacos-group/nacos-sdk-go)](https://goreportcard.com/report/github.com/nacos-group/nacos-sdk-go) ![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

---

## gancao-openapi-sdk-go

gancao-openapi-sdk-go是甘草开放平台的Go语言客户端，它实现了出入参请求的功能

# gancao_openapi
甘草开放平台sdk

## 使用限制
支持Go>=v1.23版本

## 安装
使用`go get`安装SDK：
```sh
$ go get -u github.com/zxyphp/gancao_openapi
```
## 快速使用

* GcOpenApi 初始化账户配置结构体
* 
```go
// GcOpenApi 结构体定义
client.GcOpenApi struct {
  Url       string // 开放平台地址
  AccessKey string // 账户 access key
  SecretKey string // 秘钥 secret key
}
```

* ApiRequest 请求参数结构体
```go
// GcOpenApi 结构体定义
client.ApiRequest struct {
  Package string // 请求地址 Package
  Class   string // 请求地址 Class
  Params  map[string]interface{} // 请求参数信息
}
```

* 请求示例
```go

// 初始化 GcOpenApi 配置
apiConfig := client.GcOpenApi{
AccessKey: "OU022B11I095379K7",
SecretKey: "46cb91b0f20c3e19",
}

// 创建生产环境的 API 客户端，区分环境
apiClient := client.NewGcOpenApi(apiConfig, false)

// 创建 API 请求
apiRequest := client.ApiRequest{
Package: "igc_base.example.template",
Class:   "GET_USER_INFO",
Params: map[string]interface{}{
  "name": "王小明",
  "age":  43,
  },
}

// 调用 API 示例
result, err := apiClient.ExecApi(apiRequest)
if err != nil {
fmt.Println("Error:", err)
return
}
fmt.Println("Result:", result)

```

## 文档
甘草开放平台文档 相关信息可以查看文档 [open-api doc](https://apidoc.igancao.com/home/api-standard.html).

Sdk文档可以查看 [openapi sdk](https://apidoc.igancao.com/home/api-standard.html#sdk%E4%B8%8B%E8%BD%BD).

## 联系我们
* 有问题请联系 (wechat) [zxy_coding]
* zhangxy@gancao.com

