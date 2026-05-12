package test

import (
	"blockEmulator/consensus_shard/pbft_all"
	"blockEmulator/params"
	"blockEmulator/supervisor"
	"strconv"
	"time"
)

// TEST case
//nid, _ := strconv.ParseUint(os.Args[1], 10, 64)
//nnm, _ := strconv.ParseUint(os.Args[2], 10, 64)
//sid, _ := strconv.ParseUint(os.Args[3], 10, 64)
//snm, _ := strconv.ParseUint(os.Args[4], 10, 64)
//test.TestPBFT(nid, nnm, sid, snm)

func TestPBFT(nid, nnm, sid, snm uint64) {
	params.ShardNum = int(snm)
	for i := uint64(0); i < snm; i++ {
		if _, ok := params.IPmap_nodeTable[i]; !ok {
			params.IPmap_nodeTable[i] = make(map[uint64]string)
		}
		for j := uint64(0); j < nnm; j++ {
			params.IPmap_nodeTable[i][j] = "127.0.0.1:" + strconv.Itoa(8800+int(i)*100+int(j))
		}
	}
	params.IPmap_nodeTable[params.DeciderShard] = make(map[uint64]string)
	params.IPmap_nodeTable[params.DeciderShard][0] = "127.0.0.1:18800"

	pcc := &params.ChainConfig{
		ChainID:        sid,
		NodeID:         nid,
		ShardID:        sid,
		Nodes_perShard: uint64(params.NodesInShard),
		ShardNums:      snm,
		BlockSize:      uint64(params.MaxBlockSize_global),
		BlockInterval:  uint64(params.Block_Interval),
		InjectSpeed:    uint64(params.InjectSpeed),
	}

	if nid == 12345678 {
		lsn := new(supervisor.Supervisor)
		lsn.NewSupervisor("127.0.0.1:18800", pcc, "Relay", "TPS_Relay", "Latency_Relay", "CrossTxRate_Relay", "TxNumberCount_Relay")
		time.Sleep(10000 * time.Millisecond)
		go lsn.SupervisorTxHandling()
		lsn.TcpListen()
		return
	}

	worker := pbft_all.NewPbftNode(sid, nid, pcc, "Relay")
	time.Sleep(5 * time.Second)
	if nid == 0 {
		go worker.Propose()
		worker.TcpListen()
	} else {
		worker.TcpListen()
	}
}
