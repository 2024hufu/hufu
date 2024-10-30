package handler

import (
	"hufu/controller"
	"hufu/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Transfer 处理转账请求
func Transfer(c *gin.Context) {
	var req model.Transaction
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用控制器处理转账
	from, err := controller.GetWalletByID(req.FromWalletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	to, err := controller.GetWalletByID(req.ToWalletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	tx, err := controller.TransferFunds(from, to, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tx})
}

// GetTransferHistory 获取转账历史
func GetTransferHistory(c *gin.Context) {
	var req model.Transaction
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, err := controller.GetTransferHistory(req.FromWalletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
