package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wangyaodream/gerty-goblog/pkg/config"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"github.com/wangyaodream/gerty-goblog/pkg/view"
)

type AgentController struct {
}

func agent(msg string) string {
	baseUrl := "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	apikey := config.GetString("app.apikey")
	data := map[string]interface{}{
		"model": "glm-4-flash",
		"messages": []map[string]interface{}{
			{"role": "user", "content": msg},
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
	fmt.Println(apikey)

	client := &http.Client{}
	resp, err := client.Do(res)

	if err != nil {
		logger.LogError(err)
	}

	defer resp.Body.Close()

	// 读取响应数据
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func (*AgentController) Home(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{
		"Test": "this is test",
	}, "agent.chat")
}

func (*AgentController) Chat(w http.ResponseWriter, r *http.Request) {
	msg := r.PostFormValue("body")
	result := agent(msg)
	logger.LogInfo(result)
}
