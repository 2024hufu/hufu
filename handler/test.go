package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/proxy"
)

func TestAdd(c *gin.Context) {
	// 定义请求结构体
	var requestBody struct {
		Add string `json:"add"`
	}

	// 解析 JSON 请求体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "解析 JSON 请求失败: " + err.Error(),
		})
		return
	}

	// 创建 SOCKS5 代理
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建代理失败: " + err.Error(),
		})
		return
	}

	// 创建带有代理的 HTTP 客户端
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	// 准备请求
	req, err := http.NewRequest(
		"POST",
		"http://10.77.110.184:8080/add",
		strings.NewReader(requestBody.Add), // 只发送 add 字段的值
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建请求失败: " + err.Error(),
		})
		return
	}

	// 设置 Content-Type
	req.Header.Set("Content-Type", "text/plain") // 修改为纯文本格式

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "请求失败: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "读取响应失败: " + err.Error(),
		})
		return
	}

	// 返回结果
	c.String(resp.StatusCode, string(respBody))
}
