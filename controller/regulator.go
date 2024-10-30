package controller

import (
	"fmt"

	"hufu/model"
	"hufu/supervisor"
	"hufu/utils"
)

type Regulator struct {
	maxTransactionAmount float64
}

func NewRegulator() *Regulator {
	return &Regulator{
		maxTransactionAmount: 10000,
	}
}

func (r *Regulator) SendAlert(tx *model.Transaction, evidence string) error {
	return nil
}

func (r *Regulator) CheckTransaction(tx *model.Transaction, w *model.Wallet) error {
	if tx.Amount > r.maxTransactionAmount {
		return fmt.Errorf("交易金额超过允许的最大值: %f", r.maxTransactionAmount)
	}
	return nil
}

func ProcessPrivateKey(walletID uint) ([]string, error) {
	// 获取私钥
	walletKey, err := GetWalletKeyByWalletID(walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet key: %v", err)
	}

	parts, err := utils.SharePrivateKey(walletKey.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to share private key: %v", err)
	}

	res := []string{}

	for i, part := range parts {
		encryptedPart, err := utils.EncryptData(supervisor.JuryInstance.Nodes[i].PublicKey, part)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt part %d: %v", i, err)
		}
		// 存储加密后的分片
		id := fmt.Sprintf("%d-%d", walletID, i)
		name := fmt.Sprintf("key_share_%s", id)
		err = supervisor.JuryInstance.Nodes[i].StoreEncryptedKeyShare(id, name, encryptedPart)
		res = append(res, encryptedPart)
		if err != nil {
			return nil, fmt.Errorf("failed to store encrypted part %d: %v", i, err)
		}
	}

	return res, nil
}

// GetAbnormalTransactions 获取所有异常交易
func GetAbnormalTransactions() ([]model.AbnormalTransaction, error) {
	var transactions []model.AbnormalTransaction

	// 从数据库中查询标记为异常的交易
	result := model.DB.Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}
