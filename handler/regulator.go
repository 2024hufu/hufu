package handler

import (
	"hufu/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckTransaction(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "success",
		},
	)
}

func GetPrivateKey(c *gin.Context) {
	type PrivateKeyRequest struct {
		WalletID uint   `json:"wallet_id" binding:"required"`
		Evidence string `json:"evidence" binding:"required"`
	}
	var req PrivateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}

	res, err := controller.ProcessPrivateKey(req.WalletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "处理私钥失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    res,
	})
}

// 获取异常交易
func GetAbnormalTransaction(c *gin.Context) {
	// 从数据库获取异常交易列表
	abnormalTxs, err := controller.GetAbnormalTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取异常交易失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    abnormalTxs,
	})
}
