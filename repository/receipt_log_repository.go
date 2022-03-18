package repository

import (
	"context"
	"errors"
	"github.com/ethereum-blockchain-service/model"
	"gorm.io/gorm"
	"sync"
)

var receiptLogRepo *ReceiptLogRepo
var receiptLogRepoOnce sync.Once

type ReceiptLogRepo struct {
	db *gorm.DB
}

func NewReceiptLogRepo(db *gorm.DB) ReceiptLogRepository {
	receiptLogRepoOnce.Do(func() {
		receiptLogRepo = &ReceiptLogRepo{
			db: db,
		}
	})

	return receiptLogRepo
}

func (repo *ReceiptLogRepo) Insert(ctx context.Context, transactionReceipt []*model.ReceiptLog) error {
	return repo.db.Create(transactionReceipt).Error
}

func (repo *ReceiptLogRepo) GetByTxHash(ctx context.Context, txHash []byte) ([]*model.ReceiptLog, error) {
	var result []*model.ReceiptLog
	if err := repo.db.Model(model.ReceiptLog{}).Where("tx_hash = ?", txHash).Find(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
