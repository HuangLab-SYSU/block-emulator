package params

var (
	Block_Interval      = 5000   // The time interval for generating a new block
	MaxBlockSize_global = 2000   // The maximum number of transactions a block contains
	InjectSpeed         = 2000   // The speed of transaction injection
	TotalDataSize       = 100000 // The total number of txs to be injected
	BatchSize           = 16000  // The supervisor read a batch of txs then send them. The size of a batch is 'BatchSize'
	BrokerNum           = 10
	NodesInShard        = 4
	ShardNum            = 4
	DataWrite_path      = "./result/"              // Measurement data result output path
	LogWrite_path       = "./log"                  // Log output path
	SupervisorAddr      = "127.0.0.1:18800"        // Supervisor ip address
	FileInput           = `./selectedTxs_300K.csv` // The raw BlockTransaction data path
)
