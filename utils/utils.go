package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/SSSaaS/sssa-golang"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

const (
	MINIMUM = 3 // 最小份额数
	SHARES5 = 5 // 总份额数
)

const SupervisorPrivateKey = "43878f814c6753e43c1bd91db187e9399551e50876b7d24f7aba2cc467f88458"
const SupervisorPublicKey = "066583fe9369c70280b2af181e9b6d87eb63848f7af4ac1444dcc774e11805630dfc07918bcd80803a38f77f4b6f415e1d4e2596a79ecacc83f9a0ad95645326"

const ProxyPublicKey = "066583fe9369c70280b2af181e9b6d87eb63848f7af4ac1444dcc774e11805630dfc07918bcd80803a38f77f4b6f415e1d4e2596a79ecacc83f9a0ad95645326"
const ProxyPrivateKey = "43878f814c6753e43c1bd91db187e9399551e50876b7d24f7aba2cc467f88458"

// encryptData encrypts data using ECIES with the provided ECDSA public key.
func EncryptData(key string, data string) (string, error) {
	// 将公钥字符串解码为字节数组
	publicKeyBytes, err := hex.DecodeString(key)
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

// DecryptData decrypts data using ECIES with the provided ECDSA private key.
func DecryptData(rawData string) (string, error) {
	// Decode the base64 encoded ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	// Decode the private key from hex
	privateKeyBytes, err := hex.DecodeString(ProxyPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %v", err)
	}

	// Parse the private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// Convert to ECIES private key
	eciesPrivKey := ecies.ImportECDSA(privateKey)

	// Decrypt the data
	plaintext, err := eciesPrivKey.Decrypt(ciphertext, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %v", err)
	}

	return string(plaintext), nil
}

func GenerateHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// generateKeys 生成随机密钥对
func GenerateKeys() (privateKeyString, publicKeyString, address string) {
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

	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return hexutil.Encode(privateKeyBytes)[2:], hexutil.Encode(publicKeyBytes)[4:], address
}

// SharePrivateKey Shamir's secret sharing https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing
func SharePrivateKey(PrivateKey string) ([]string, error) {
	result, err := sssa.Create(MINIMUM, SHARES5, PrivateKey)
	return result, err
}

// func ReadConfig() *client.Config {
// 	privateKey, _ := hex.DecodeString("145e247e170ba3afd6ae97e88f00dbc976c2345d511b0f6713355d19d8b80b58")
// 	// disable ssl of node rpc
// 	config := &client.Config{
// 		IsSMCrypto:  false,
// 		GroupID:     "group0",
// 		DisableSsl:  false,
// 		PrivateKey:  privateKey,
// 		Host:        "127.0.0.1",
// 		Port:        20200,
// 		TLSCaFile:   "./conf/ca.crt",
// 		TLSKeyFile:  "./conf/sdk.key",
// 		TLSCertFile: "./conf/sdk.crt",
// 	}

// 	return config
// }
