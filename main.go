package main

import (
	"github.com/FISCO-BCOS/go-sdk/hufu/controller"
	"github.com/FISCO-BCOS/go-sdk/hufu/model"
)

func main() {
	model.SetupDB()
	w1 := controller.GetWalletByName("wallet01")
	w2 := controller.GetWalletByName("wallet02")
	err := controller.TransferFunds(w1.Address, w2.Address, 2)
	if err != nil {
		return
	}
}
