package model

type ReceiptLog struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	TxHash      []byte `json:"tx_hash"`
	BlockNumber int64  `json:"block_number"`
	LogIndex    int64  `json:"log_index"`
	Data        []byte `json:"data"`
}
