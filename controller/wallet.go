package controller

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/FISCO-BCOS/go-sdk/hufu/model"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

// generateKeys 生成随机密钥对
func generateKeys() ([]byte, []byte, string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("private key: ", hexutil.Encode(privateKeyBytes)[2:]) // privateKey in hex without "0x"

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("publick key: ", hexutil.Encode(publicKeyBytes)[4:]) // publicKey in hex without "0x"

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("address: ", address) // account address
	return privateKeyBytes, publicKeyBytes, address
}

func NewWallet(name string, balance float64) (*model.Wallet, error) {
	privateKey, publicKey, address := generateKeys()

	//config := &client.Config{
	//	IsSMCrypto:  false,
	//	GroupID:     "group0",
	//	PrivateKey:  privateKey,
	//	Host:        "127.0.0.1",
	//	Port:        20200,
	//	TLSCaFile:   "../conf/ca.crt",
	//	TLSKeyFile:  "../conf/sdk.key",
	//	TLSCertFile: "../conf/sdk.crt",
	//}
	//client, err := client.DialContext(context.Background(), config)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(client)

	w := &model.Wallet{
		Name:       name,
		Address:    address,
		PublicKey:  hexutil.Encode(privateKey)[2:],
		PrivateKey: hexutil.Encode(publicKey)[4:],
		Balance:    balance,
	}

	// Name cannot be repeated
	if GetWalletByName(name).Name != "" {
		return nil, fmt.Errorf("wallet name already exists: %s", name)
	}

	model.DB.Create(w)

	return w, nil
}

func GetWalletByName(Name string) *model.Wallet {
	var w model.Wallet
	model.DB.Where("name = ?", Name).First(&w)
	return &w
}

func GetWalletByAddress(Address string) *model.Wallet {
	var w model.Wallet
	model.DB.Where("address = ?", Address).First(&w)
	return &w
}

func SetWalletBalance(Name string, balance float64) {
	model.DB.Model(&model.Wallet{}).Where("name = ?", Name).Update("balance", balance)
}
