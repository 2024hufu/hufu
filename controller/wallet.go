package controller

import (
	"fmt"
	"hufu/model"
	"hufu/utils"
)

func NewWallet(name string, balance float64) (*model.Wallet, error) {
	w := &model.Wallet{
		Name:    name,
		Balance: balance,
	}

	// Name cannot be repeated
	// if _, err := GetWalletByName(name); err == nil {
	// 	return nil, fmt.Errorf("wallet name already exists: %s", name)
	// }
	// if _, err := GetWalletByID(w.ID); err == nil {
	// 	return nil, fmt.Errorf("wallet id already exists: %d", w.ID)
	// }

	if err := model.DB.Create(w).Error; err != nil {
		return nil, err
	}

	privateKey, publicKey, _ := utils.GenerateKeys()
	walletKey := &model.WalletKey{
		WalletID:   w.ID,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	if err := model.DB.Create(walletKey).Error; err != nil {
		return nil, err
	}

	return w, nil
}

func GetWalletByName(Name string) (*model.Wallet, error) {
	var w model.Wallet
	if err := model.DB.Where("name = ?", Name).First(&w).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func GetWalletByID(ID uint) (*model.Wallet, error) {
	var w model.Wallet
	if err := model.DB.Where("id = ?", ID).First(&w).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func SetWalletBalance(ID uint, balance float64) error {
	return model.DB.Model(&model.Wallet{}).Where("id = ?", ID).Update("balance", balance).Error
}

func GetWalletKeyByWalletID(walletID uint) (*model.WalletKey, error) {
	var wk model.WalletKey
	if err := model.DB.Where("wallet_id = ?", walletID).First(&wk).Error; err != nil {
		return nil, err
	}
	return &wk, nil
}

// UpdateWallet 更新钱包信息
func UpdateWallet(walletID uint, name string, balance float64) error {
	// 首先检查钱包是否存在
	wallet, err := GetWalletByID(walletID)
	if err != nil {
		return fmt.Errorf("wallet not found: %v", err)
	}

	// 更新钱包信息
	wallet.Name = name
	wallet.Balance = balance

	// 保存更新
	if err := model.DB.Save(wallet).Error; err != nil {
		return fmt.Errorf("update wallet failed: %v", err)
	}

	return nil
}
