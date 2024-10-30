// package supervisor

// import (
// 	"context"
// 	"hufu/contract/KeyShare"
// 	"hufu/utils"
// 	"log"

// 	"github.com/FISCO-BCOS/go-sdk/client"
// 	"github.com/ethereum/go-ethereum/common"
// )

// const (
// 	Contract1 = "0x1234567890abcdef"
// 	Contract2 = "0xabcdef1234567890"
// 	Contract3 = "0x7890abcdef1234567"
// )

// func CallContract1(evidence string) error {
// 	config := utils.ReadConfig()
// 	client, err := client.DialContext(context.Background(), config)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	contractAddress := common.HexToAddress("0x6849f21d1e455e9f0712b1e99fa4fcd23758e8f1") // 0x481D3A1dcD72cD618Ea768b3FbF69D78B46995b0

// 	instance, err := KeyShare.NewKeyShare(contractAddress, client)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// call the contract function
// 	session := &KeyShare.KeyShareSession{
// 		Contract:     instance,
// 		CallOpts:     *client.GetCallOpts(),
// 		TransactOpts: *client.GetTransactOpts(),
// 	}

// 	//session.Insert("1", "key1", "value1")
// 	//
// 	//session.Insert("2", "key2", "value2")
// 	//session.Insert("3", "key3", "value3")

// 	v1, v2, v3 := session.Select("2")
// 	log.Println(v1, v2, v3)

// 	return nil
// }

// func SetPrivateKey(name, privateKey string) error {
// 	return nil
// }

// func GetPrivateKey(name string) (string, error) {
// 	return "", nil
// }
