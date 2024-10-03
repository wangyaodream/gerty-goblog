package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wangyaodream/gerty-goblog/pkg/auth"
	"github.com/wangyaodream/gerty-goblog/pkg/config"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
)

type AgentController struct {
}

func (ac AgentController) Agent(w http.ResponseWriter, r *http.Request) {
	baseUrl := "https://open.bigmodel.cn/api/paas/v4/chat/completions'"
	currentUser := auth.User()
	apikey := config.GetString("apikey")
	data := map[string]interface{}{
		"model": "glm-4-flash",
		"messages": []map[string]interface{}{
			{"role": "user", "content": "你好，请介绍下自己"},
		},
	}
	jsonData, _ := json.Marshal(data)

	// 发送api请求并获取内容
	res, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("请求失败!")
	}
	res.Header.Set("Content-Type", "application/json")
	res.Header.Set("Authorization", "Bearer "+apikey)

	client := &http.Client{}
	resp, err := client.Do(res)

	if err != nil {
		logger.LogError(err)
	}

	defer resp.Body.Close()

	// 读取响应数据
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Hello! %v result:\n", currentUser.Name)
	fmt.Println(body)
}
