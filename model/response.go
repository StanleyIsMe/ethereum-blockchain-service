package model

type Resp struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Blocks struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	BlockNum   int64  `json:"block_num"`
	BlockHash  []byte `json:"block_hash"`
	BlockTime  int64  `json:"block_time"`
	ParentHash []byte `json:"parent_hash"`
}

type BlockWithTxHash struct {
	BlockNum   int64    `json:"block_num"`
	BlockHash  []byte   `json:"block_hash"`
	BlockTime  int64    `json:"block_time"`
	ParentHash []byte   `json:"parent_hash"`
	TxHash     []string `json:"transactions"`
}

type TxWithLogs struct {
	TxHash []byte     `json:"tx_hash"`
	From   []byte     `json:"from"`
	To     []byte     `json:"to"`
	Nonce  int64      `json:"nonce"`
	Data   []byte     `json:"data"`
	Value  string     `json:"value"`
	Logs   []EventLog `json:"logs"`
}

type EventLog struct {
	Index int64  `json:"index"`
	Data  []byte `json:"data"`
}
