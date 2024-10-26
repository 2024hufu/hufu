package controller

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hufu/errors"
	"hufu/model"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

const (
	TransactionStatusFailed  = "failed"
	TransactionStatusSuccess = "success"
)

const ProxyWalletCount = 3

func TransferFunds(from *model.Wallet, to *model.Wallet, amount float64) error {
	transaction := &model.Transaction{
		TransactionHash: generateHash(fmt.Sprintf("%s-%s-%f-%d", from.Address, to.Address, amount, time.Now().Unix())),
		SenderAddress:   from.Address,
		ReceiverAddress: to.Address,
		Amount:          amount,
		Timestamp:       time.Now().Unix(),
		Status:          TransactionStatusFailed,
	}
	encryptedTransaction, err := createEncryptedTransaction(transaction, from)
	if err != nil {
		return err
	}
	desensitizedTransaction := createDesensitizedTransaction(transaction)

	model.DB.Create(encryptedTransaction)
	model.DB.Create(desensitizedTransaction)

	return ProxyTransferFunds(from, to, amount)
}

func ProxyTransferFunds(from *model.Wallet, to *model.Wallet, amount float64) error {
	// 检查余额
	if from.Balance < amount {
		return errors.ErrInsufficientBalance
	}

	// 获取随机代理钱包
	proxyWallets := GlobalWalletPool.GetRandomWallets(ProxyWalletCount)
	if len(proxyWallets) < ProxyWalletCount {
		return fmt.Errorf("not enough proxies, only %d proxies available", len(proxyWallets))
	}

	// 计算每个代理钱包需要转账的金额
	amountPerProxy := amount / float64(ProxyWalletCount)

	// 开始事务
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 从发送方转账到代理钱包
	for _, proxyWallet := range proxyWallets {
		if err := transferBetweenWallets(from, proxyWallet, amountPerProxy); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 从代理钱包转账到接收方
	for _, proxyWallet := range proxyWallets {
		if err := transferBetweenWallets(proxyWallet, to, amountPerProxy); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func transferBetweenWallets(from *model.Wallet, to *model.Wallet, amount float64) error {
	transaction := &model.Transaction{
		TransactionHash: generateHash(fmt.Sprintf("%s-%s-%f-%d", from.Address, to.Address, amount, time.Now().Unix())),
		SenderAddress:   from.Address,
		ReceiverAddress: to.Address,
		Amount:          amount,
		Timestamp:       time.Now().Unix(),
		Status:          TransactionStatusFailed,
	}

	if from == nil || to == nil {
		return errors.ErrWalletNotFound
	}

	if from.Balance < amount {
		return errors.ErrInsufficientBalance
	}

	// 更新钱包余额
	from.Balance -= amount
	to.Balance += amount

	// 更新钱包余额
	if err := model.DB.Model(from).Where("address = ?", from.Address).Update("balance", from.Balance).Error; err != nil {
		return err
	}
	if err := model.DB.Model(to).Where("address = ?", to.Address).Update("balance", to.Balance).Error; err != nil {
		return err
	}
	transaction.Status = TransactionStatusSuccess
	// 更新真实交易表
	if err := model.DB.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

func createEncryptedTransaction(t *model.Transaction, w *model.Wallet) (*model.EncryptedTransaction, error) {
	// 使用公钥加密交易信息
	encryptedTransactionHash, err := encryptData(t.TransactionHash, w.PublicKey)
	if err != nil {
		return nil, err
	}
	encryptedSenderAddress, err := encryptData(t.SenderAddress, w.PublicKey)
	if err != nil {
		return nil, err
	}
	encryptedReceiverAddress, err := encryptData(t.ReceiverAddress, w.PublicKey)
	if err != nil {
		return nil, err
	}
	encryptedAmount, err := encryptData(fmt.Sprintf("%f", t.Amount), w.PublicKey)
	if err != nil {
		return nil, err
	}
	encryptedTimestamp, err := encryptData(fmt.Sprintf("%d", t.Timestamp), w.PublicKey)
	if err != nil {
		return nil, err
	}

	return &model.EncryptedTransaction{
		TransactionHash: encryptedTransactionHash,
		SenderAddress:   encryptedSenderAddress,
		ReceiverAddress: encryptedReceiverAddress,
		Amount:          encryptedAmount,
		Timestamp:       encryptedTimestamp,
	}, nil
}

// TODO: 脱敏金额和时间范围需要根据实际情况进行调整
func createDesensitizedTransaction(t *model.Transaction) *model.DesensitizedTransaction {
	return &model.DesensitizedTransaction{
		TransactionHash: t.TransactionHash,
		SenderHash:      generateHash(t.SenderAddress),
		ReceiverHash:    generateHash(t.ReceiverAddress),
		AmountRange:     fmt.Sprintf("%f-%f", t.Amount-10, t.Amount+10),
		TimeRange:       fmt.Sprintf("%d-%d", t.Timestamp-10, t.Timestamp+10),
	}
}

func generateHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// encryptData encrypts data using ECIES with the provided ECDSA public key.
func encryptData(data string, publicKeyString string) (string, error) {
	// 将公钥字符串解码为字节数组
	publicKeyBytes, err := hex.DecodeString(publicKeyString)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %v", err)
	}

	// 直接使用解码后的公钥字节
	// 加上0x04前缀，表示未压缩的公钥
	publicKeyBytes = append([]byte{0x04}, publicKeyBytes...)
	publicKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal public key: %v", err)
	}

	// Convert to ECIES public key
	eciesPubKey := ecies.ImportECDSAPublic(publicKey)

	// Encrypt the data
	ciphertext, err := ecies.Encrypt(rand.Reader, eciesPubKey, []byte(data), nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %v", err)
	}

	// Encode the ciphertext to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
