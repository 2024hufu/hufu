package model

// Transaction 交易
type Transaction struct {
	TransactionHash string  `json:"transaction_hash"`
	SenderAddress   string  `json:"sender_address"`
	ReceiverAddress string  `json:"receiver_address"`
	Amount          float64 `json:"amount"`
	Timestamp       int64   `json:"timestamp"`
	Status          string  `json:"status"`
}

// EncryptedTransaction 加密交易
type EncryptedTransaction struct {
	TransactionHash string `json:"transaction_hash"`
	SenderAddress   string `json:"sender_address"`
	ReceiverAddress string `json:"receiver_address"`
	Amount          string `json:"amount"`
	Timestamp       string `json:"timestamp"`
}

// DesensitizedTransaction 脱敏交易
type DesensitizedTransaction struct {
	TransactionHash string `json:"transaction_hash"`
	SenderHash      string `json:"sender_hash"`
	ReceiverHash    string `json:"receiver_hash"`
	AmountRange     string `json:"amount_range"`
	TimeRange       string `json:"time_range"`
}
