package controller

import (
	"fmt"
	"github.com/ethereum-blockchain-service/model"
	"github.com/ethereum-blockchain-service/service/api/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BlockController struct {
	blockUseCase *usecase.BlockUseCase
}

func NewBlockController() *BlockController {
	return &BlockController{
		blockUseCase: usecase.NewBlockUseCase(),
	}
}

func (ctrl *BlockController) GetBlocksByLimit(c *gin.Context) {
	ctx := c.Request.Context()
	limitString, ok := c.GetQuery("limit")
	if !ok {
		c.JSON(http.StatusBadRequest, model.Resp{
			Message: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	limit, _ := strconv.Atoi(limitString)
	fmt.Println(limit, "!")
	result, err := ctrl.blockUseCase.FindBlockByLatest(ctx, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Resp{
			Message: http.StatusText(http.StatusNotFound),
		})
		return
	}

	response := model.Blocks{Blocks: make([]model.Block, len(result))}
	for index, block := range result {
		response.Blocks[index] = model.Block{
			BlockNum:   block.Number,
			BlockHash:  block.Hash,
			BlockTime:  block.Time,
			ParentHash: block.ParentHash,
		}
	}

	c.JSON(http.StatusOK, model.Resp{
		Message: http.StatusText(http.StatusOK),
		Data:    response,
	})
}

func (ctrl *BlockController) GetBlocksById(c *gin.Context) {
	ctx := c.Request.Context()
	blockNumber, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Resp{
			Message: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	header, transactions, err := ctrl.blockUseCase.FindBlockAndTransactionByNumber(ctx, int64(blockNumber))
	if err != nil || header == nil {
		c.JSON(http.StatusNotFound, model.Resp{
			Message: http.StatusText(http.StatusNotFound),
		})
		return
	}

	response := model.BlockWithTxHash{
		BlockNum:   header.Number,
		BlockHash:  header.Hash,
		BlockTime:  header.Time,
		ParentHash: header.ParentHash,
		TxHash:     make([]string, len(transactions)),
	}

	for index, transaction := range transactions {
		response.TxHash[index] = string(transaction.Hash)
	}

	c.JSON(http.StatusOK, model.Resp{
		Message: http.StatusText(http.StatusOK),
		Data:    response,
	})
}

func (ctrl *BlockController) GetTransactionByTxHash(c *gin.Context) {
	ctx := c.Request.Context()
	txHash := []byte(c.Param("txHash"))
	transaction, eventLogs, err := ctrl.blockUseCase.FindTransactionAndLogByTxHash(ctx, txHash)
	if err != nil || transaction == nil {
		c.JSON(http.StatusNotFound, model.Resp{
			Message: http.StatusText(http.StatusNotFound),
		})
		return
	}

	response := model.TxWithLogs{
		TxHash: txHash,
		From:   transaction.From,
		To:     transaction.To,
		Nonce:  transaction.Nonce,
		Data:   transaction.Data,
		Value:  transaction.Value,
		Logs:   make([]model.EventLog, len(eventLogs)),
	}
	for index, eventLog := range eventLogs {
		response.Logs[index] = model.EventLog{
			Index: eventLog.LogIndex,
			Data:  eventLog.Data,
		}
	}

	c.JSON(http.StatusOK, model.Resp{
		Message: http.StatusText(http.StatusOK),
		Data:    response,
	})
}
