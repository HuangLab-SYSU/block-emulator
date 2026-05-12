package params

var (
	Block_Interval      = 5000   // generate new block interval
	MaxBlockSize_global = 2000   // the block contains the maximum number of transactions
	InjectSpeed         = 5000   // the transaction inject speed
	TotalDataSize       = 500000 // the total number of txs
	BatchSize           = 5000   // supervisor read a batch of txs then send them, it should be larger than inject speed
	BrokerNum           = 50
	NodesInShard        = 4
	ShardNum            = 4
	IterNum_B2E         = 5
	Brokerage           = 0.1
	DataWrite_path      = "./result/"                                                                 // measurement data result output path
	LogWrite_path       = "./log"                                                                     // log output path
	SupervisorAddr      = "127.0.0.1:18800"                                                           //supervisor ip address
	FileInput           = `./data/selectedTxs_300K.csv` //the raw BlockTransaction data path
)
