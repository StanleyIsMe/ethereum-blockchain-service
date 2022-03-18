package usecase

import (
	"context"
	"fmt"
	"github.com/ethereum-blockchain-service/model"
	"github.com/ethereum-blockchain-service/repository"
	"sync"
)

var blockUseCase *BlockUseCase
var blockUseCaseInstanceOnce sync.Once

type BlockUseCase struct {
	repository *repository.Repository
}

func NewBlockUseCase() *BlockUseCase {
	blockUseCaseInstanceOnce.Do(func() {
		blockUseCase = &BlockUseCase{
			repository: repository.NewRepository(),
		}
	})
	return blockUseCase
}

func (use *BlockUseCase) FindBlockByLatest(ctx context.Context, num int) ([]*model.BlockHeader, error) {
	return use.repository.BlockHeader.GetLatest(ctx, num)
}

func (use *BlockUseCase) FindBlockAndTransactionByNumber(ctx context.Context, number int64) (*model.BlockHeader, []*model.Transactions, error) {
	header, err := use.repository.BlockHeader.GetByNumber(ctx, number)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	transactions, err := use.repository.Transactions.GetByNumber(ctx, header.Number)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return header, transactions, nil
}

func (use *BlockUseCase) FindTransactionAndLogByTxHash(ctx context.Context, txHash []byte) (*model.Transactions, []*model.ReceiptLog, error) {
	transaction, err := use.repository.Transactions.GetByTxHash(ctx, txHash)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	eventLogs, err := use.repository.ReceiptLog.GetByTxHash(ctx, txHash)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	return transaction, eventLogs, nil
}
