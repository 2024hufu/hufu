package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hufu/model"
	"time"
)

func TransferFunds(from string, to string, amount float64) error {
	fromWallet := GetWalletByAddress(from)
	toWallet := GetWalletByAddress(to)

	transaction := &model.Transaction{
		TransactionHash: generateTransactionHash(from, to, amount),
		SenderAddress:   from,
		ReceiverAddress: to,
		Amount:          amount,
		Timestamp:       time.Now().Unix(),
		Status:          "failed",
	}

	if fromWallet == nil {
		model.DB.Create(transaction)
		return fmt.Errorf("sender wallet not found: %s", from)
	}
	if toWallet == nil {
		model.DB.Create(transaction)
		return fmt.Errorf("receiver wallet not found: %s", to)
	}

	if fromWallet.Balance < amount {
		model.DB.Create(transaction)
		return fmt.Errorf("insufficient balance in wallet: %s", from)
	}

	fromWallet.Balance -= amount
	toWallet.Balance += amount

	model.DB.Where("address = ?", fromWallet.Address).Updates(fromWallet)
	model.DB.Where("address = ?", toWallet.Address).Updates(toWallet)

	transaction.Status = "success"
	model.DB.Create(transaction)

	// Generate sensitive data
	sensitiveData := generateSensitiveData(transaction)
	CreateSensitiveData(sensitiveData)

	return nil
}

func generateTransactionHash(sender string, receiver string, amount float64) string {
	data := fmt.Sprintf("%s-%s-%f-%d", sender, receiver, amount, time.Now().Unix())

	hash := sha256.New()
	hash.Write([]byte(data))

	return hex.EncodeToString(hash.Sum(nil))
}

func generateSensitiveData(tx *model.Transaction) *model.SensitiveData {
	return &model.SensitiveData{
		TransactionHash: tx.TransactionHash,
		SenderHash:      "generateHash(tx.SenderAddress)",
		ReceiverHash:    "generateHash(tx.ReceiverAddress)",
		AmountRange:     "0",
		TimeRange:       "0",
	}
}
