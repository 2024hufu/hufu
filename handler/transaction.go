package handler

import (
	"encoding/json"
	"hufu/config"
	"hufu/controller"
	"hufu/model"
	"hufu/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Req struct {
	WalletID   uint   `json:"wallet_id"`
	PrivateKey string `json:"private_key"`
}

// EncryptedTransfer 处理加密转账请求
func EncryptedTransfer(c *gin.Context) {
	// 获取加密请求数据
	var encryptedReq struct {
		EncryptedData string `json:"encrypted_data"`
		KeyID         string `json:"key_id"`
	}

	if err := c.ShouldBindJSON(&encryptedReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 使用私钥解密数据
	privateKey := config.GlobalConfig.Tee.PrivateKey
	decrypted, err := utils.RSADecrypt(encryptedReq.EncryptedData, privateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "解密失败: " + err.Error(),
		})
		return
	}

	// 解析解密后的数据
	var txData struct {
		FromWalletID uint    `json:"from_wallet_id"`
		ToWalletID   uint    `json:"to_wallet_id"`
		Amount       float64 `json:"amount"`
	}

	if err := json.Unmarshal([]byte(decrypted), &txData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "解析数据失败: " + err.Error(),
		})
		return
	}
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
