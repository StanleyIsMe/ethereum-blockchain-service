package model

type Transactions struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Hash        []byte `json:"hash"`
	BlockHash   []byte `json:"block_hash"`
	From        []byte `json:"from"`
	To          []byte `json:"to"`
	Nonce       int64  `json:"nonce"`
	BlockNumber int64  `json:"block_number"`
	Data        []byte `json:"data"`
	Value       string `json:"value"`
}
