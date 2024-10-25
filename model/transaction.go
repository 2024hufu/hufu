package model

type Transaction struct {
	TransactionHash string  `json:"transaction_hash"`
	SenderAddress   string  `json:"sender_address"`
	ReceiverAddress string  `json:"receiver_address"`
	Amount          float64 `json:"amount"`
	Timestamp       int64   `json:"timestamp"`
	Status          string  `json:"status"`
}
