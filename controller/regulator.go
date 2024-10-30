package controller

import (
	"fmt"
	"log"

	"hufu/model"
	"hufu/utils"
)

var juryPrivateKeys = []string{
	"2e0d5bb53302e8157bcf66e8c200d289321f70e87b0c3ef86c7f7b32a24cde31",
	"e82ad9284cb0e28fadde1baded3ad07686cba2798e7e47288ee3ba5eb8402573",
	"a2f4b2644a5c543c25147813fcc6ecd38ea0da24e9315e861b153b117db09fe2",
	"ad06b8d8ec4a693b6c51bf7747f45c8b740eb71bb58867052c311aa9e171a7b2",
	"0db45ccfa10b887ac0a72f667b94f1cd00ca78375fe6cbc5055db7590840336c",
}

var juryPublicKeys = []string{
	"02fcc4219340ca21e121ae234c0f55701944cd3e4b64ddfff72f6e5ddd55c12b161d628af00c80f165e53af66499dfa38fa4eef76edab6c485c43a6ec35831eb",
	"1080675212fd218e0597a74d4b05ce39e2011441ca4b556c752990754e52fde2df9f91af59d2cee19c00dc96ee49709d90847be001d4818615922561ed3d230e",
	"50fd454198b70992bd6bf2d8560a8df660fbce66577d6548ccc1d732ceb38bbd1138b72b34d1bca884e6bb9739cda5db4c44f6d425b89218a13325a8aaa4b41f",
	"1935c942890992bc6f018c4b0d6436512f79b7981cecaf69db1c6eb28f3d3d4744668fff78d26c7144af5c39ba53706db6f6a785bb936a6b6cddfa35e533a686",
	"d549fb67d5c9a1dc63a29f102f6cc67714da1ac4171035250a8ba749c97438e92e497de53a918b403da8217eab23f8a3d4dde499314e41461f120f6b1795db01",
}

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
		encryptedPart, err := utils.EncryptData(juryPublicKeys[i], part)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt part %d: %v", i, err)
		}
		// 存储加密后的分片
		err = storeEncryptedPart(encryptedPart, i)
		if err != nil {
			return nil, fmt.Errorf("failed to store encrypted part %d: %v", i, err)
		}
	}

	return res, nil
}

// storeEncryptedPart 存储加密后的分片
func storeEncryptedPart(encryptedPart string, partIndex int) error {
	// TODO: 实现具体存储逻辑存储在区块链上
	log.Println("k:", partIndex, "v:", encryptedPart)
	return nil
}
