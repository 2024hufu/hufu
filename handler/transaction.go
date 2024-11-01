package handler

import (
	"hufu/controller"
	"hufu/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Req struct {
	WalletID   uint   `json:"wallet_id"`
	PrivateKey string `json:"private_key"`
}

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
	var req struct {
		WalletID uint `json:"wallet_id" binding:"required"`
		Page     int  `json:"page"`
		PageSize int  `json:"page_size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	result, err := controller.GetTransferHistory(req.WalletID, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
}

// GetEncryptedTransaction 获取加密交易信息
func GetEncryptedTransaction(c *gin.Context) {
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedTx, err := controller.GetEncryptedTransaction(req.WalletID, req.PrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": encryptedTx,
	})
}

// GetDesensitizedTransaction 获取脱敏交易记录
func GetDesensitizedTransaction(c *gin.Context) {
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, err := controller.GetDesensitizedTransaction(req.WalletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

// GetReceivedTransactions 获取收款记录
func GetReceivedTransactions(c *gin.Context) {
	var req struct {
		WalletID uint `json:"wallet_id" binding:"required"`
		Page     int  `json:"page"`
		PageSize int  `json:"page_size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	result, err := controller.GetReceivedTransactions(req.WalletID, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
}

// GetTransactionStats 获取交易统计信息
func GetTransactionStats(c *gin.Context) {
	var req struct {
		WalletID uint `json:"wallet_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	stats, err := controller.GetTransactionStats(req.WalletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "获取统计信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}
