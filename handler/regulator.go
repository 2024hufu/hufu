package handler

import (
	"hufu/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrivateKeyRequest struct {
	WalletID uint `json:"wallet_id" binding:"required"`
}

func CheckTransaction(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "success",
		},
	)
}

func GetPrivateKey(c *gin.Context) {
	var req PrivateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}

	_, err := controller.ProcessPrivateKey(req.WalletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "处理私钥失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "私钥提交成功",
	})
}
