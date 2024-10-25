package model

type Wallet struct {
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	PrivateKey string  `json:"private_key"`
	PublicKey  string  `json:"public_key"`
	Balance    float64 `json:"balance"`
}
