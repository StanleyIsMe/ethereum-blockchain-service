package repository

import (
	"context"
	"errors"
	"github.com/ethereum-blockchain-service/model"
	"sync"

	"gorm.io/gorm"
)

var blockHeaderRepo *BlockHeaderRepo
var blockHeaderRepoOnce sync.Once

type BlockHeaderRepo struct {
	db *gorm.DB
}

func NewBlockHeader(db *gorm.DB) BlockHeaderRepository {
	blockHeaderRepoOnce.Do(func() {
		blockHeaderRepo = &BlockHeaderRepo{
			db: db,
		}
	})

	return blockHeaderRepo
}

func (repo *BlockHeaderRepo) Create(ctx context.Context, header *model.BlockHeader) error {
	return repo.db.Create(header).Error
}

func (repo *BlockHeaderRepo) GetLatest(ctx context.Context, num int) ([]*model.BlockHeader, error) {
	var result []*model.BlockHeader
	if err := repo.db.Model(model.BlockHeader{}).Order("number desc").Limit(num).Find(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (repo *BlockHeaderRepo) GetByNumber(ctx context.Context, number int64) (*model.BlockHeader, error) {
	var result *model.BlockHeader
	if err := repo.db.Model(model.BlockHeader{}).Where("number = ?", number).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
