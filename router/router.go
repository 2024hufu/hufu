package router

import (
	"hufu/handler"

	"github.com/gin-gonic/gin"
)

func InitHufuRouter(r *gin.Engine) {
	hufu := r.Group("/api/v1/hufu")
	{
		// 钱包相关路由
		wallet := hufu.Group("/wallet")
		{
			wallet.POST("/create", handler.CreateWallet) // 创建钱包
			wallet.GET("/", handler.GetWallet)           // 获取单个钱包
			wallet.POST("/update", handler.UpdateWallet) // 更新钱包
		}

		// 转账相关路由
		tx := hufu.Group("/tx")
		{
			tx.POST("/transfer", handler.Transfer)                // 转账
			tx.GET("/history", handler.GetTransferHistory)        // 获取转账历史
			tx.GET("/encrypted", handler.GetEncryptedTransaction) // 获取加密交易
		}
	}
}

func InitRegulatorRouter(r *gin.Engine) {
	regulator := r.Group("/api/v1/regulator")
	{
		regulator.POST("/alert", handler.CheckTransaction)         // 检查交易
		regulator.GET("/private-key", handler.GetPrivateKey)       // 获取私钥
		regulator.GET("/abnormal", handler.GetAbnormalTransaction) // 获取异常交易
	}
}
