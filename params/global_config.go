package params

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var (
	// The following parameters can be set in main.go.
	// default values:
	NodesInShard = 4 // # of Nodes in a shard.
	ShardNum     = 4 // # of shards.
)

var (
	ConsensusMethod = 0 // ConsensusMethod an Integer, which indicates the choice ID of methods / consensuses. Value range: [0, 4), representing [CLPA_Broker, CLPA, Broker, Relay]"

	Block_Interval      = 5000   // The time interval for generating a new block
	MaxBlockSize_global = 2000   // The maximum number of transactions a block contains
	InjectSpeed         = 2000   // The speed of transaction injection
	TotalDataSize       = 160000 // The total number of txs to be injected
	TxBatchSize         = 16000  // The supervisor read a batch of txs then send them. The size of a batch is 'TxBatchSize'
	BrokerNum           = 10     // The # of Broker accounts used in Broker / CLPA_Broker.

	ExpDataRootDir     = "expTest"                     // The root dir where the experimental data should locate.
	DataWrite_path     = ExpDataRootDir + "/result/"   // Measurement data result output path
	LogWrite_path      = ExpDataRootDir + "/log"       // Log output path
	DatabaseWrite_path = ExpDataRootDir + "/database/" // database write path

	SupervisorAddr = "127.0.0.1:18800"        // Supervisor ip address
	DatasetFile    = `./selectedTxs_300K.csv` // The raw BlockTransaction data path

	ReconfigTimeGap = 50 // The time gap between epochs. This variable is only used in CLPA / CLPA_Broker now.
)

// read from file
type globalConfig struct {
	ConsensusMethod int `json:"ConsensusMethod"`

	ExpDataRootDir string `json:"ExpDataRootDir"`

	BlockInterval      int `json:"Block_Interval"`
	MaxBlockSizeGlobal int `json:"BlockSize"`
	InjectSpeed        int `json:"InjectSpeed"`
	TotalDataSize      int `json:"TotalDataSize"`

	TxBatchSize     int    `json:"TxBatchSize"`
	BrokerNum       int    `json:"BrokerNum"`
	DatasetFile     string `json:"DatasetFile"`
	ReconfigTimeGap int    `json:"ReconfigTimeGap"`
}

func ReadConfigFile() {
	// read configurations from paramsConfig.json
	data, err := os.ReadFile("paramsConfig.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	var config globalConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// output configurations
	fmt.Printf("Config: %+v\n", config)

	// set configurations to params
	ConsensusMethod = config.ConsensusMethod

	ExpDataRootDir = config.ExpDataRootDir
	DataWrite_path = ExpDataRootDir + "/result/"
	LogWrite_path = ExpDataRootDir + "/log"
	DatabaseWrite_path = ExpDataRootDir + "/database/"

	Block_Interval = config.BlockInterval
	MaxBlockSize_global = config.MaxBlockSizeGlobal

	InjectSpeed = config.InjectSpeed
	TotalDataSize = config.TotalDataSize
	TxBatchSize = config.TxBatchSize

	BrokerNum = config.BrokerNum
	DatasetFile = config.DatasetFile

	ReconfigTimeGap = config.ReconfigTimeGap
}
