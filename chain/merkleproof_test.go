package chain

import (
	"blockEmulator/core"
	"blockEmulator/params"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"
)

func TestMerkleProof(t *testing.T) {
	txProof, txHash := buildBlockChain()
	clearBlockchainData()
	// fmt.Printf("%v", txProof)
	if ok, err := TxProofVerify(txHash, &txProof); !ok {
		log.Panic("Fail to verify ", err.Error())
	}
}

func buildBlockChain() (TxProofResult, []byte) {
	// build a blockchain
	fp := params.DatabaseWrite_path + "mptDB/ldb/s0/N0"
	db, err := rawdb.NewLevelDBDatabase(fp, 0, 1, "accountState", false)
	if err != nil {
		log.Panic(err)
	}
	params.ShardNum = 1
	pcc := &params.ChainConfig{
		ChainID:        0,
		NodeID:         0,
		ShardID:        0,
		Nodes_perShard: 1,
		ShardNums:      1,
		BlockSize:      uint64(params.MaxBlockSize_global),
		BlockInterval:  uint64(params.Block_Interval),
		InjectSpeed:    uint64(params.InjectSpeed),
	}
	CurChain, _ := NewBlockChain(pcc, db)
	CurChain.PrintBlockChain()

	// update blokchain
	TxForProof := core.NewTransaction("00000000001", "00000000002", big.NewInt(1234), 126526, time.Now())
	for i := 0; i < 4; i++ {
		// add a special tx for further proof validation.
		if i == 2 {
			CurChain.Txpool.AddTx2Pool(TxForProof)
		}
		// add txs
		for j := 0; j < 1000; j++ {
			CurChain.Txpool.AddTx2Pool(core.NewTransaction("00000000001", "00000000002", big.NewInt(1000), uint64(i*1234+j+213), time.Now()))
		}
		b := CurChain.GenerateBlock(int32(i))
		CurChain.AddBlock(b)
	}

	// get proof of this Tx
	txProofResult := CurChain.TxProofGenerate(TxForProof.TxHash)

	// close blockchain
	CurChain.CloseBlockChain()

	return txProofResult, TxForProof.TxHash
}

func clearBlockchainData() {
	// clear test data file
	err := os.RemoveAll(params.ExpDataRootDir)
	if err != nil {
		fmt.Printf("Failed to delete directory: %v\n", err)
		return
	}
}
