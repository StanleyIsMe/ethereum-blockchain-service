package main

import (
	"context"
	"fmt"
	"github.com/ethereum-blockchain-service/model"
	"github.com/ethereum-blockchain-service/repository"
	"github.com/ethereum-blockchain-service/tools/db"
	"github.com/ethereum-blockchain-service/tools/routine"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	url := "https://data-seed-prebsc-2-s3.binance.org:8545"
	c, err := ethclient.Dial(url)
	if err == nil {
		log.Info("Connect to eth client successfully", "url", url)
	} else {
		log.Warn("Failed to dial eth client", "url", url, "err", err)
	}

	db.Init()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		WaitShutdown(cancel)
	}()
	indexer := NewIndexer(c)
	if err := indexer.Listen(ctx, 23027); err != nil {
		log.Error(err.Error())
	}

	cancel()
}

func WaitShutdown(cancelFunc context.CancelFunc) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh)
	for s := range signalCh {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM:
			log.Info("Signal: %v, Begin Shutdown", s)
			go TimeOut()
			close(signalCh)
			cancelFunc()
			log.Info("Shutdown Finished")
			os.Exit(0)
		default:
		}
	}
}

// 超時強制關閉
func TimeOut() {
	for {
		select {
		case <-time.After(time.Second * 60):
			log.Info("Shutdown Timeout!!!!!!!!")
			os.Exit(0)
		}
	}
}

func NewIndexer(client *ethclient.Client) *Indexer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Indexer{
		client:        client,
		currentHeader: nil,
		ctx:           ctx,
		cancel:        cancel,
		repository:    repository.NewRepository(),
	}
}

type Indexer struct {
	ctx           context.Context
	cancel        context.CancelFunc
	client        *ethclient.Client
	currentHeader *model.BlockHeader
	repository    *repository.Repository
}

func (idx *Indexer) Listen(ctx context.Context, fromBlock int64) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	blockReceiveChannel := make(chan *types.Block, 1000)
	go idx.ReceiveBlock(ctx, blockReceiveChannel)
	if err := idx.LoadLocalBlock(ctx, fromBlock, blockReceiveChannel); err != nil {
		fmt.Println(err)
		return err
	}

	listenCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	subHeaderCh := make(chan *types.Header, 100)
	workerPool := routine.NewWorkerPool(20)
	defer workerPool.Shutdown()

	//sub, err := idx.client.SubscribeNewHead(ctx, subHeaderCh)
	//if err != nil {
	//	fmt.Println("fail SubscribeNewHead")
	//	return err
	//}

	wg := sync.WaitGroup{}
	for {
		select {
		case header := <-subHeaderCh:
			if idx.currentHeader.Number >= header.Number.Int64() {
				fmt.Println("ignore old header")
				continue
			}

			currentNumber := idx.currentHeader.Number
			for currentNumber < header.Number.Int64() {
				wg.Add(1)
				workerPool.AddTask(func(params []interface{}) {
					defer wg.Done()
					number := big.NewInt(params[0].(int64))
					block, err := idx.client.BlockByNumber(ctx, number)
					if err != nil {
						log.Error(err.Error())
						return
					}
					blockReceiveChannel <- block
				}, currentNumber)
				currentNumber++
			}
			wg.Wait()
			idx.currentHeader = model.Header(header)
		//case err := <-sub.Err():
		//	return err
		case <-listenCtx.Done():
			return listenCtx.Err()
		default:
			block, err := idx.client.BlockByNumber(ctx, nil)
			if err != nil {
				log.Error("BlockByNumber error: %v", err)
			}
			subHeaderCh <- block.Header()
			time.Sleep(5 * time.Second)
		}

	}
}

func (idx *Indexer) ReceiveBlock(ctx context.Context, blockChannel <-chan *types.Block) error {
	receiveCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case block, ok := <-blockChannel:
			if !ok {
				return nil
			}

			log.Info("Receive Block: ", block.Number(), block.Hash())
			header := block.Header()
			err := idx.repository.BlockHeader.Create(ctx, model.Header(header))
			if err != nil {
				log.Error("create header Error : %v", err)
				continue
			}

			for _, transaction := range block.Transactions() {
				transactions, receiptLogs, err := idx.ParseTransaction(ctx, transaction, header)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				if err := idx.repository.Transactions.Create(ctx, transactions); err != nil {
					log.Error(err.Error())
					continue
				}

				if err := idx.repository.ReceiptLog.Insert(ctx, receiptLogs); err != nil {
					log.Error(err.Error())
					continue
				}
			}
		case <-receiveCtx.Done():
			return receiveCtx.Err()
		}
	}

}

func (idx *Indexer) ParseTransaction(ctx context.Context, tx *types.Transaction, header *types.Header) (*model.Transactions, []*model.ReceiptLog, error) {
	receipt, err := idx.client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		log.Error("TransactionReceipt error: %v", err)
		return nil, nil, err
	}

	if receipt != nil {
		signer := types.MakeSigner(params.MainnetChainConfig, header.Number)
		msg, err := tx.AsMessage(signer, nil)
		if err != nil {
			log.Error(err.Error())
			return nil, nil, err
		}

		// Transaction
		var to []byte
		if msg.To() != nil {
			to = msg.To().Bytes()
		}

		value := ""
		if tx.Value() != nil {
			value = tx.Value().String()
		}
		transaction := &model.Transactions{
			Hash:        tx.Hash().Bytes(),
			BlockHash:   header.Hash().Bytes(),
			From:        msg.From().Bytes(),
			To:          to,
			Nonce:       int64(tx.Nonce()),
			BlockNumber: header.Number.Int64(),
			Data:        tx.Data(),
			Value:       value,
		}

		receiptLogs := make([]*model.ReceiptLog, len(receipt.Logs))
		for index, txLog := range receipt.Logs {
			receiptLogs[index] = &model.ReceiptLog{
				TxHash:      tx.Hash().Bytes(),
				BlockNumber: header.Number.Int64(),
				LogIndex:    int64(txLog.Index),
				Data:        txLog.Data,
			}
		}

		return transaction, receiptLogs, nil
	}
	return nil, nil, nil
}

func (idx *Indexer) LoadLocalBlock(ctx context.Context, fromBlock int64, receiveCh chan<- *types.Block) error {
	headerList, err := idx.repository.BlockHeader.GetLatest(ctx, 1)
	if err != nil {
		return err
	}

	if len(headerList) == 0 {
		block, err := idx.client.BlockByNumber(ctx, big.NewInt(fromBlock))
		if err != nil {
			return err
		}
		receiveCh <- block
		idx.currentHeader = model.Header(block.Header())
		return nil
	}

	idx.currentHeader = headerList[0]
	return nil
}
