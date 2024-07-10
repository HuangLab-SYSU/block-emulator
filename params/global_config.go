package params

var (
	// The following parameters can be set in main.go.
	// default values:
	NodesInShard   = 4         // # of Nodes in a shard.
	ShardNum       = 4         // # of shards.
	ExpDataRootDir = "expTest" // The root dir where the experimental data should locate.
)

var (
	Block_Interval      = 5000   // The time interval for generating a new block
	MaxBlockSize_global = 2000   // The maximum number of transactions a block contains
	InjectSpeed         = 2000   // The speed of transaction injection
	TotalDataSize       = 160000 // The total number of txs to be injected
	BatchSize           = 16000  // The supervisor read a batch of txs then send them. The size of a batch is 'BatchSize'
	BrokerNum           = 10     // The # of Broker accounts used in Broker / CLPA_Broker.

	DataWrite_path     = ExpDataRootDir + "/result/"   // Measurement data result output path
	LogWrite_path      = ExpDataRootDir + "/log"       // Log output path
	DatabaseWrite_path = ExpDataRootDir + "/database/" // database write path

	SupervisorAddr = "127.0.0.1:18800"        // Supervisor ip address
	FileInput      = `./selectedTxs_300K.csv` // The raw BlockTransaction data path

	ReconfigTimeGap = 50 // The time gap between epochs. This variable is only used in CLPA / CLPA_Broker now.
)
