package handler

import (
	"hufu/controller"
	"hufu/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateWallet 创建钱包
func CreateWallet(c *gin.Context) {
	var req model.Wallet
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := controller.NewWallet(req.Name, req.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": wallet})
}

// GetWallet 获取单个钱包详情
func GetWallet(c *gin.Context) {
	var req model.Wallet
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := controller.GetWalletByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": wallet})
}

// UpdateWallet 更新钱包信息
func UpdateWallet(c *gin.Context) {
	var req model.Wallet

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 调用 controller 层更新钱包
	err := controller.UpdateWallet(req.ID, req.Name, req.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "更新钱包失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新钱包成功",
	})
}
