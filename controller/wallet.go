package controller

import (
	"crypto/ecdsa"
	"fmt"
	"hufu/model"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewWallet(name string, balance float64) (*model.Wallet, error) {
	privateKey, publicKey, address := generateKeys()

	w := &model.Wallet{
		Name:       name,
		Address:    address,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Balance:    balance,
	}

	// Name cannot be repeated
	if _, err := GetWalletByName(name); err == nil {
		return nil, fmt.Errorf("wallet name already exists: %s", name)
	}

	if err := model.DB.Create(w).Error; err != nil {
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

func GetWalletByAddress(Address string) (*model.Wallet, error) {
	var w model.Wallet
	if err := model.DB.Where("address = ?", Address).First(&w).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func SetWalletBalance(Name string, balance float64) error {
	return model.DB.Model(&model.Wallet{}).Where("name = ?", Name).Update("balance", balance).Error
}

// generateKeys 生成随机密钥对
func generateKeys() (string, string, string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return hexutil.Encode(privateKeyBytes)[2:], hexutil.Encode(publicKeyBytes)[4:], address
}
