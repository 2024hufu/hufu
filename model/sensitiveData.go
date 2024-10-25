package model

type SensitiveData struct {
	TransactionHash string `json:"transaction_hash"`
	SenderHash      string `json:"sender_hash"`
	ReceiverHash    string `json:"receiver_hash"`
	AmountRange     string `json:"amount_range"`
	TimeRange       string `json:"time_range"`
}
