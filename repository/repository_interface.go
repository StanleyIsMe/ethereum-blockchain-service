package repository

import (
	"context"
	"github.com/ethereum-blockchain-service/model"
	"github.com/ethereum-blockchain-service/tools/db"
	"gorm.io/gorm"
	"sync"
)

type (
	BlockHeaderRepository interface {
		Create(ctx context.Context, header *model.BlockHeader) error
		GetLatest(ctx context.Context, num int) ([]*model.BlockHeader, error)
		GetByNumber(ctx context.Context, number int64) (*model.BlockHeader, error)
	}

	TransactionsRepository interface {
		Create(ctx context.Context, transactions *model.Transactions) error
		GetByTxHash(ctx context.Context, txHash []byte) (*model.Transactions, error)
		GetByNumber(ctx context.Context, number int64) ([]*model.Transactions, error)
	}

	ReceiptLogRepository interface {
		Insert(ctx context.Context, receiptLog []*model.ReceiptLog) error
		GetByTxHash(ctx context.Context, txHash []byte) ([]*model.ReceiptLog, error)
	}
)

var repositoryOnce sync.Once
var repositoryObj *Repository

type Repository struct {
	db *gorm.DB

	BlockHeader  BlockHeaderRepository
	Transactions TransactionsRepository
	ReceiptLog   ReceiptLogRepository
}

func NewRepository() *Repository {
	repositoryOnce.Do(func() {
		repositoryObj = &Repository{
			db:           db.DB.Instance,
			BlockHeader:  NewBlockHeader(db.DB.Instance),
			Transactions: NewTransactionsRepo(db.DB.Instance),
			ReceiptLog:   NewReceiptLogRepo(db.DB.Instance),
		}
	})
	return repositoryObj
}
