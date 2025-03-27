package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func getUserinfo() (string, string, string, string) {
	// 打开JSON文件
	file, err := os.Open("../datas/api.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// 读取文件内容
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// 定义结构体以匹配JSON结构
	type APIConfig struct {
		APIKey     string `json:"api_key"`
		SecretKey  string `json:"secret_key"`
		Passphrase string `json:"passphrase"`
	}

	var config APIConfig

	// 解析JSON内容
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	// flag是实盘与模拟盘的切换参数
	// flag = "1" // 模拟盘 demo trading
	flag := "0" // 实盘 real trading

	return config.APIKey, config.SecretKey, config.Passphrase, flag
}

func main() {
	apiKey, secretKey, passphrase, flag := getUserinfo()
	fmt.Println("API Key:", apiKey)
	fmt.Println("Secret Key:", secretKey)
	fmt.Println("Passphrase:", passphrase)
	fmt.Println("Flag:", flag)
}
