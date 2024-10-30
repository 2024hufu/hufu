package router

import (
	"hufu/handler"

	"github.com/gin-gonic/gin"
)

func InitRegulatorRouter(r *gin.Engine) {
	regulator := r.Group("/api/v1/regulator")
	{
		regulator.POST("/alert", handler.CheckTransaction)   // 检查交易
		regulator.GET("/private-key", handler.GetPrivateKey) // 获取私钥
	}
}
