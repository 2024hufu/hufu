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
			tx.POST("/transfer", handler.Transfer)         // 转账
			tx.GET("/history", handler.GetTransferHistory) // 获取转账历史
		}
	}
}
