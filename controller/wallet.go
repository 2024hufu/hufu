package controller

import (
	"fmt"
	"hufu/model"
	"hufu/utils"
	"time"
)

type WalletStats struct {
	TodayTransactions int64 `json:"today_transactions"` // 今日交易次数
	TotalTransactions int64 `json:"total_transactions"` // 总交易次数
}

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

func GetWalletStats(walletID uint) (*WalletStats, error) {
	// 获取今天的开始时间（零点）
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var todayCount int64
	if err := model.DB.Model(&model.Transaction{}).
		Where("(from_wallet_id = ? OR to_wallet_id = ?) AND created_at >= ? AND type = ?",
			walletID, walletID, todayStart, model.DirectTransaction).
		Count(&todayCount).Error; err != nil {
		return nil, err
	}

	var totalCount int64
	if err := model.DB.Model(&model.Transaction{}).
		Where("(from_wallet_id = ? OR to_wallet_id = ?) AND type = ?",
			walletID, walletID, model.DirectTransaction).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	return &WalletStats{
		TodayTransactions: todayCount,
		TotalTransactions: totalCount,
	}, nil
}
