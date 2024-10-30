package supervisor

import (
	"fmt"
	"hufu/model"

	"github.com/gin-gonic/gin"
)

func HandleAlert(c *gin.Context) error {
	var alertData struct {
		Transaction *model.Transaction `json:"transaction"`
		Evidence    string             `json:"evidence"`
	}

	// 从请求体中读取JSON数据
	if err := c.ShouldBindJSON(&alertData); err != nil {
		return err
	}

	fmt.Println("收到警报: ", alertData)

	// request private key
	// RequestPrivateKey(alertData.Transaction.SenderAddress)

	return nil
}
