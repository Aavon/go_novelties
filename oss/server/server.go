package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Aavon/go_novelties/oss"
)

func main() {
	http.HandleFunc("/authorize", Authorize)
	log.Println("start...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("listen: %v", err)
	}
}

var (
	ossConfig = oss.OssOptions{
		Endpoint:         "oss-cn-beijing.aliyuncs.com",
		EndpointInternal: "oss-cn-beijing-internal.aliyuncs.com",
		AccessKeyId:      "",
		AccessKeySecret:  "",
		BucketName:       "",
	}
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	ossKey := fmt.Sprintf("test/%d", time.Now().Unix())
	policyToken := oss.GetPolicyToken(ossConfig, ossKey, "", time.Now().Unix()+300)
	policyStr, err := json.Marshal(policyToken)
	if err != nil {
		log.Fatal(err)
	}
	// 测试
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(policyStr)
	if err != nil {
		log.Println(err)
	}
}
