package handler

import (
	"hufu/config"
	"hufu/utils"

	"github.com/gin-gonic/gin"
)

// GetEncryptionKeys 获取加密密钥的处理函数
func GetEncryptionKeys(c *gin.Context) {
	c.JSON(200, gin.H{
		"public_key": config.GlobalConfig.Tee.PublicKey,
	})
}

func EncryptData(c *gin.Context) {
	type EncryptDataRequest struct {
		Data string `json:"data"`
	}

	var request EncryptDataRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := utils.EncryptData(request.Data, config.GlobalConfig.Tee.PublicKey)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": res})
}
