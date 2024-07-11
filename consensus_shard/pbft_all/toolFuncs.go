package pbft_all

import (
	"blockEmulator/core"
	"blockEmulator/message"
	"blockEmulator/params"
	"blockEmulator/shard"
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// set 2d map, only for pbft maps, if the first parameter is true, then set the cntPrepareConfirm map,
// otherwise, cntCommitConfirm map will be set
func (p *PbftConsensusNode) set2DMap(isPrePareConfirm bool, key string, val *shard.Node) {
	if isPrePareConfirm {
		if _, ok := p.cntPrepareConfirm[key]; !ok {
			p.cntPrepareConfirm[key] = make(map[*shard.Node]bool)
		}
		p.cntPrepareConfirm[key][val] = true
	} else {
		if _, ok := p.cntCommitConfirm[key]; !ok {
			p.cntCommitConfirm[key] = make(map[*shard.Node]bool)
		}
		p.cntCommitConfirm[key][val] = true
	}
}

// get neighbor nodes in a shard
func (p *PbftConsensusNode) getNeighborNodes() []string {
	receiverNodes := make([]string, 0)
	for _, ip := range p.ip_nodeTable[p.ShardID] {
		receiverNodes = append(receiverNodes, ip)
	}
	return receiverNodes
}

func (p *PbftConsensusNode) writeCSVline(metricName []string, metricVal []string) {
	// Construct directory path
	dirpath := params.DataWrite_path + "pbft_shardNum=" + strconv.Itoa(int(p.pbftChainConfig.ShardNums))
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}

	// Construct target file path
	targetPath := fmt.Sprintf("%s/Shard%d%d.csv", dirpath, p.ShardID, p.pbftChainConfig.ShardNums)

	// Open file, create if it does not exist
	file, err := os.OpenFile(targetPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)

	// Write header if the file is newly created
	fileInfo, err := file.Stat()
	if err != nil {
		log.Panic(err)
	}
	if fileInfo.Size() == 0 {
		if err := writer.Write(metricName); err != nil {
			log.Panic(err)
		}
		writer.Flush()
	}

	// Write data
	if err := writer.Write(metricVal); err != nil {
		log.Panic(err)
	}
	writer.Flush()
}

// get the digest of request
func getDigest(r *message.Request) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

// calculate TCL
func computeTCL(txs []*core.Transaction, commitTS time.Time) int64 {
	ret := int64(0)
	for _, tx := range txs {
		ret += commitTS.Sub(tx.Time).Milliseconds()
	}
	return ret
}
