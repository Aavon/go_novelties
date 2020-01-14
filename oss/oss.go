package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"time"
)

/*
	参考 aliyun-oss-appserver-js-master
	集成和使用更加方便
	生成签名所需的信息可通过官方库(github.com/aliyun/aliyun-oss-go-sdk/oss)的Bucket获取
	接口文档参考： https://help.aliyun.com/document_detail/31988.html?spm=a2c4g.11186623.2.12.67c718eeX5SwTT#reference-smp-nsw-wdb
*/

type OssConfig struct {
	Name    string     `json:"name" mapstructure:"name"`
	Options OssOptions `json:"options" mapstructure:"options"`
}

type OssOptions struct {
	Endpoint         string `json:"endpoint" mapstructure:"endpoint"`
	EndpointInternal string `json:"endpoint_internal" mapstructure:"endpoint_internal"`
	AccessKeyId      string `json:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret  string `json:"access_key_secret" mapstructure:"access_key_secret"`
	BucketName       string `json:"bucket_name" mapstructure:"bucket_name"`
}

type PolicyConfig struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

// 前端将以下添加到form表单中，通过POST方式上传
type PolicyToken struct {
	AccessKeyId string `json:"OSSAccessKeyId"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	ObjectKey   string `json:"key"`
	// 文件上传回调(可选)
	Callback string `json:"callback,omitempty"`
	// file
}

func GetPolicyToken(options OssOptions, objectKey, callbackUrl string, expireAt int64) PolicyToken {
	// 前端添加http/https
	host := fmt.Sprintf("%s.%s", options.BucketName, options.Endpoint)

	policyStr := getPolicyStr(expireAt, objectKey)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(options.AccessKeySecret))
	_, _ = io.WriteString(h, policyStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackBase64 string
	if callbackUrl != "" {
		var callbackParam CallbackParam
		callbackParam.CallbackUrl = callbackUrl
		callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
		callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
		callback_str, err := json.Marshal(callbackParam)
		if err != nil {
			fmt.Println("callback json err:", err)
		}
		callbackBase64 = base64.StdEncoding.EncodeToString(callback_str)
	}

	var policyToken PolicyToken
	policyToken.AccessKeyId = options.AccessKeyId
	policyToken.Host = host
	policyToken.Expire = expireAt
	policyToken.Signature = string(signedStr)
	policyToken.ObjectKey = objectKey
	policyToken.Policy = string(policyStr)
	policyToken.Callback = callbackBase64
	return policyToken
}

func getPolicyStr(expireAt int64, objectKey string) string {
	tokenExpire := time.Unix(expireAt, 0).Format("2006-01-02T15:04:05Z")

	//create post policy json
	var config PolicyConfig
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "eq")
	condition = append(condition, "$key")
	condition = append(condition, objectKey)
	config.Conditions = append(config.Conditions, condition)

	// base64
	result, _ := json.Marshal(config)
	return base64.StdEncoding.EncodeToString(result)
}
