package main

import (
	"hufu/controller"
	"hufu/model"
	"hufu/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	model.SetupDB()
	controller.InitWalletPool()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.InitHufuRouter(r)
	router.InitRegulatorRouter(r)
	r.Run(":3338")

	// w1, _ := controller.NewWallet("Alice1", 100)
	// w2, _ := controller.NewWallet("Alice2", 100)

	// w1, _ := controller.GetWalletByName("Alice1")
	// w2, _ := controller.GetWalletByName("Alice2")

	// tx, err := controller.TransferFunds(w1, w2, 10)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// txs, _ := controller.GetTransferHistory(11)
	// for _, tx := range txs {
	// 	fmt.Printf("%+v\n", tx)
	// }
}
