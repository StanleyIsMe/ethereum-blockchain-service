package model

import (
	"github.com/ethereum/go-ethereum/core/types"
	"time"
)

type BlockHeader struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	Hash       []byte    `json:"hash"`
	ParentHash []byte    `json:"parent_hash"`
	Root       []byte    `json:"root"`
	TxHash     []byte    `json:"tx_hash"`
	Number     int64     `json:"number"`
	Time       int64     `json:"time"`
	nonce      []byte    `json:"nonce"`
	CreatedAt  time.Time `json:"created_at"`
}

func Header(header *types.Header) *BlockHeader {
	return &BlockHeader{
		Hash:       header.Hash().Bytes(),
		ParentHash: header.ParentHash.Bytes(),
		Root:       header.Root.Bytes(),
		TxHash:     header.TxHash.Bytes(),
		Number:     header.Number.Int64(),
		Time:       int64(header.Time),
	}
}
