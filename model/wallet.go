package model

import "gorm.io/gorm"

// Wallet 钱包
type Wallet struct {
	gorm.Model
	Name    string  `json:"name" gorm:"type:varchar(100);not null"`
	Balance float64 `json:"balance" gorm:"type:decimal(20,8);default:0"`
}

// 钱包公私钥
type WalletKey struct {
	gorm.Model
	WalletID   uint   `json:"wallet_id" gorm:"type:int;not null"`
	PublicKey  string `json:"public_key" gorm:"type:varchar(256);not null"`
	PrivateKey string `json:"private_key" gorm:"type:varchar(256);not null"`
}
