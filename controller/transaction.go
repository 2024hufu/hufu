package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hufu/errors"
	"hufu/model"
	"hufu/utils"
	"log"
	"net/http"

	"gorm.io/gorm"
)

const (
	TransactionStatusFailed  = "failed"
	TransactionStatusSuccess = "success"
)

const ProxyWalletCount = 3

// TransferFunds 处理转账
func TransferFunds(from *model.Wallet, to *model.Wallet, amount float64) (*model.Transaction, error) {
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建并保存原始交易
	originalTx, err := createOriginalTransaction(tx, from, to, amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. 验证交易合规性
	if err := validateTransaction(tx, from, originalTx); err != nil {
		// 如果是不合规的交易，直接提交事务并返回
		// 因为validateTransaction已经更新了交易状态并创建了异常记录
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		return originalTx, nil
	}

	// 3. 创建关联交易记录
	if err := createAssociatedTransactions(tx, originalTx); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 4. 更新钱包余额
	if err := updateWalletBalances(tx, from, to, amount); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 5. 完成交易
	if err := finalizeTransaction(tx, originalTx); err != nil {
		tx.Rollback()
		return nil, err
	}

	return originalTx, tx.Commit().Error
}

// GetTransferHistory 获取转账历史
func GetTransferHistory(walletID uint) ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	// 查询与指定钱包ID相关的交易记录
	if err := model.DB.Where("from_wallet_id = ? OR to_wallet_id = ?", walletID, walletID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	var res []*model.Transaction
	for _, tx := range transactions {
		if tx.Type == model.DirectTransaction {
			res = append(res, tx)
		}
	}

	return res, nil
}

// createOriginalTransaction 创建原始交易记录
func createOriginalTransaction(tx *gorm.DB, from *model.Wallet, to *model.Wallet, amount float64) (*model.Transaction, error) {
	originalTx := &model.Transaction{
		FromWalletID: from.ID,
		ToWalletID:   to.ID,
		Amount:       amount,
		Type:         model.DirectTransaction,
		Status:       "pending",
	}
	return originalTx, tx.Create(originalTx).Error
}

// validateTransaction 验证交易合规性
func validateTransaction(tx *gorm.DB, from *model.Wallet, originalTx *model.Transaction) error {
	if originalTx.Amount > 10000 {
		// 使用传入的事务对象创建异常交易记录
		abnormal := &model.AbnormalTransaction{
			WalletID:      from.ID,
			TransactionID: originalTx.ID,
		}

		if err := tx.Create(abnormal).Error; err != nil {
			return err
		}

		// 更新原始交易状态为失败
		originalTx.Status = TransactionStatusFailed
		if err := tx.Save(originalTx).Error; err != nil {
			return err
		}

		// 异步发送告警
		go func() {
			jsonData, err := json.Marshal(abnormal)
			if err != nil {
				log.Printf("marshal abnormal transaction failed: %v", err)
				return
			}

			resp, err := http.Post("http://127.0.0.1:3338/api/v1/regulator/alert", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("send alert failed: %v", err)
				return
			}
			defer resp.Body.Close()
		}()

		return errors.ErrTransactionAmountTooLarge
	}
	return nil
}

// createAssociatedTransactions 创建关联交易记录
func createAssociatedTransactions(tx *gorm.DB, originalTx *model.Transaction) error {
	// 创建加密交易记录
	encryptedTx, err := createEncryptedTransaction(originalTx)
	if err != nil {
		return err
	}
	if err := tx.Create(encryptedTx).Error; err != nil {
		return err
	}

	// 创建脱敏交易记录
	desensitizedTx := createDesensitizedTransaction(originalTx)
	if err := tx.Create(desensitizedTx).Error; err != nil {
		return err
	}

	// 创建代交易记录
	proxyTxs, err := createProxyTransactions(originalTx, originalTx.Amount)
	if err != nil {
		return err
	}
	for _, proxyTx := range proxyTxs {
		if err := tx.Create(proxyTx).Error; err != nil {
			return err
		}
	}
	return nil
}

// updateWalletBalances 更新钱包余额
func updateWalletBalances(tx *gorm.DB, from *model.Wallet, to *model.Wallet, amount float64) error {
	if from.Balance < amount {
		return errors.ErrInsufficientBalance
	}

	from.Balance -= amount
	to.Balance += amount

	if err := tx.Save(from).Error; err != nil {
		return err
	}
	return tx.Save(to).Error
}

// finalizeTransaction 完成交易
func finalizeTransaction(tx *gorm.DB, originalTx *model.Transaction) error {
	originalTx.Status = "completed"
	return tx.Save(originalTx).Error
}

// createProxyTransactions 创建代理交易记录
func createProxyTransactions(originalTx *model.Transaction, amount float64) ([]*model.Transaction, error) {
	// 获取代理钱包
	proxyWallets, err := GlobalWalletPool.GetRandomWallets(ProxyWalletCount)
	if err != nil {
		return nil, err
	}

	// 拆分金额
	amounts := splitAmount(amount)
	if len(amounts) != ProxyWalletCount {
		return nil, fmt.Errorf("金额拆分数量与代理钱包数量不匹配")
	}

	proxyTxs := make([]*model.Transaction, 0, ProxyWalletCount*2)

	// 1. 创建从原始钱包到代理钱包的交易
	for i, proxyWallet := range proxyWallets {
		toProxyTx := &model.Transaction{
			FromWalletID: originalTx.FromWalletID,
			ToWalletID:   proxyWallet.ID,
			Amount:       amounts[i],
			Type:         model.ToProxyTransaction,
			Status:       "completed",
		}
		proxyTxs = append(proxyTxs, toProxyTx)

		// 2. 创建从代理钱包到目标钱包的交易
		fromProxyTx := &model.Transaction{
			FromWalletID: proxyWallet.ID,
			ToWalletID:   originalTx.ToWalletID,
			Amount:       amounts[i],
			Type:         model.FromProxyTransaction,
			Status:       "completed",
		}
		proxyTxs = append(proxyTxs, fromProxyTx)
	}

	return proxyTxs, nil
}

// splitAmount 将金额拆分成多个小额
func splitAmount(amount float64) []float64 {
	// 这里实现金额拆分的逻辑
	// 示例：简单地将金额平均拆分为3份
	part := amount / 3
	return []float64{part, part, amount - 2*part}
}

func createEncryptedTransaction(t *model.Transaction) (*model.EncryptedTransaction, error) {
	// 从数据库读取公钥和私钥
	walletKey, err := GetWalletKeyByWalletID(t.FromWalletID)
	if err != nil {
		return nil, err
	}

	// 加密交易
	encryptedFromWalletID, err := utils.EncryptData(walletKey.PublicKey, fmt.Sprintf("%d", t.FromWalletID))
	if err != nil {
		return nil, err
	}

	encryptedToWalletID, err := utils.EncryptData(walletKey.PublicKey, fmt.Sprintf("%d", t.ToWalletID))
	if err != nil {
		return nil, err
	}

	encryptedAmount, err := utils.EncryptData(walletKey.PublicKey, fmt.Sprintf("%f", t.Amount))
	if err != nil {
		return nil, err
	}

	return &model.EncryptedTransaction{
		EncryptedFromWalletID: encryptedFromWalletID,
		EncryptedToWalletID:   encryptedToWalletID,
		EncryptedAmount:       encryptedAmount,
	}, nil
}

// TODO: 脱敏金额和时间范围需要根据实际情况进行调整
func createDesensitizedTransaction(t *model.Transaction) *model.DesensitizedTransaction {
	return &model.DesensitizedTransaction{
		FromWalletID: t.FromWalletID,
		ToWalletID:   t.ToWalletID,
		AmountRange:  fmt.Sprintf("%f-%f", t.Amount-10, t.Amount+10),
		TimeRange:    fmt.Sprintf("%d-%d", t.CreatedAt.Unix()-10, t.CreatedAt.Unix()+10),
	}
}

// GetEncryptedTransaction 获取加密交易信息
func GetEncryptedTransaction(walletID uint, privateKey string) ([]*model.Transaction, error) {
	walletKey, err := GetWalletKeyByWalletID(walletID)
	if err != nil {
		return nil, err
	}

	if walletKey.PrivateKey != privateKey {
		return nil, errors.ErrPrivateKeyInvalid
	}

	// 从数据库读取
	var txs []*model.Transaction
	if err := model.DB.Where("from_wallet_id = ? OR to_wallet_id = ?", walletID, walletID).Find(&txs).Error; err != nil {
		return nil, err
	}

	res := make([]*model.Transaction, 0)
	for _, tx := range txs {
		if tx.Type == model.DirectTransaction {
			res = append(res, tx)
		}
	}

	return res, nil
}
