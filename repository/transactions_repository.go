package repository

import (
	"context"
	"errors"
	"github.com/ethereum-blockchain-service/model"
	"gorm.io/gorm"
	"sync"
)

var transactionsRepo *TransactionsRepo
var transactionsRepoOnce sync.Once

type TransactionsRepo struct {
	db *gorm.DB
}

func NewTransactionsRepo(db *gorm.DB) TransactionsRepository {
	transactionsRepoOnce.Do(func() {
		transactionsRepo = &TransactionsRepo{
			db: db,
		}
	})

	return transactionsRepo
}

func (repo *TransactionsRepo) Create(ctx context.Context, transactions *model.Transactions) error {
	return repo.db.Create(transactions).Error
}

func (repo *TransactionsRepo) GetByTxHash(ctx context.Context, txHash []byte) (*model.Transactions, error) {
	var result *model.Transactions
	if err := repo.db.Model(model.Transactions{}).Where("hash = ?", txHash).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (repo *TransactionsRepo) GetByNumber(ctx context.Context, number int64) ([]*model.Transactions, error) {
	var result []*model.Transactions
	if err := repo.db.Model(model.Transactions{}).Where("block_number = ?", number).Find(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
