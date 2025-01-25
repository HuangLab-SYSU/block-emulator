package chain

import (
	"blockEmulator/core"
	"blockEmulator/params"
	"bytes"
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
	if ok, err := TxProofVerify(txHash, &txProof); !ok && err != nil {
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
			txToAdd := core.NewTransaction("00000000001", "00000000002", big.NewInt(1000), uint64(i*1234+j+213), time.Now())
			if bytes.Equal(txToAdd.TxHash, TxForProof.TxHash) {
				log.Panic(fmt.Errorf("conflict hash"))
			}
			CurChain.Txpool.AddTx2Pool(txToAdd)
		}
		b := CurChain.GenerateBlock(int32(0))
		CurChain.AddBlock(b)
		CurChain.PrintBlockChain()
	}

	// get proof of this Tx
	txProofResult := CurChain.TxProofGenerate(TxForProof.TxHash)
	if !txProofResult.Found {
		log.Panic("cannot block found")
	}
	fmt.Printf("txProofResult: %v\n", txProofResult.BlockHeight)

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
