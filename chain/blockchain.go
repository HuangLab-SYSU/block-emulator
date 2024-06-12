package chain

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"blockEmulator/account"
	"blockEmulator/core"
	"blockEmulator/params"
	"blockEmulator/storage"

	// "blockEmulator/trie"
	"blockEmulator/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

var (
	blocktimelog, queuelenlog *csv.Writer
)

type BlockChain struct {
	ChainConfig *params.ChainConfig // Chain configuration

	CurrentBlock *core.Block // Current head of the block chain

	Storage *storage.Storage
	db      ethdb.Database // the leveldb database to store in the disk, for status trie
	Triedb  *trie.Database // the trie database which helps to store the status trie

	// StatusTrie *trie.Trie

	Tx_pool *core.Tx_pool

	TXmig1_pool *core.TXmig1_pool

	TXmig2_pool *core.TXmig2_pool

	TXann_pool *core.TXann_pool

	TXns_pool *core.TXns_pool
}

func NewBlockChain(chainConfig *params.ChainConfig) (*BlockChain, error) {
	var err error
	if chainConfig.NodeID == "N0" {
		csvFile, err := os.Create("./log/" + chainConfig.ShardID + "_blocktime.csv")
		if err != nil {
			log.Panic(err)
		}
		// defer csvFile.Close()
		blocktimelog = csv.NewWriter(csvFile)
		blocktimelog.Write([]string{"height", "block", "tx", "mig1", "mig2", "ann", "ns"})
		blocktimelog.Flush()

		csvFile1, err := os.Create("./log/" + chainConfig.ShardID + "_queueLen.csv")
		if err != nil {
			log.Panic(err)
		}
		queuelenlog = csv.NewWriter(csvFile1)
		queuelenlog.Write([]string{"block", "queueLen"})
		queuelenlog.Flush()
	}

	fmt.Printf("%v\n", chainConfig)
	bc := &BlockChain{
		ChainConfig: chainConfig,
		Storage:     storage.NewStorage(chainConfig),
		Tx_pool:     core.NewTxPool(),

		//不停止时要用的
		TXmig1_pool: core.NewTXmig1Pool(),
		TXmig2_pool: core.NewTXmig2Pool(),
		TXann_pool:  core.NewTXannPool(),
		TXns_pool:   core.NewTXnsPool(),
	}
	//filepath
	fp := "./record/triedb/" + chainConfig.ShardID + "_" + chainConfig.NodeID
	// cache=0 和 handles=1 只要小于16，都会被设为16
	bc.db, err = rawdb.NewLevelDBDatabase(fp, 0, 1, "accountState", false)
	if err != nil {
		log.Panic("cannot get the level db")
	}

	blockHash, err := bc.Storage.GetNewestBlockHash()
	if err != nil {
		if err.Error() == "newestBlockHash is not found" {
			genesisBlock := bc.NewGenesisBlock()
			bc.AddGenesisBlock(genesisBlock)
			return bc, nil
		}
		log.Panic()
	}

	// there is a blockchain in the storage
	block, err := bc.Storage.GetBlock(blockHash)
	if err != nil {
		log.Panic()
	}
	bc.CurrentBlock = block
	// stateTree, err := bc.Storage.GetStatusTree()
	// if err != nil {
	// 	log.Panic()
	// }
	// bc.StatusTrie = stateTree
	triedb := trie.NewDatabaseWithConfig(bc.db, &trie.Config{
		Cache:     0,
		Preimages: true,
	})
	bc.Triedb = triedb
	// check the existence of the trie database
	_, err = trie.New(trie.TrieID(common.BytesToHash(block.Header.StateRoot)), triedb)
	if err != nil {
		log.Panic()
	}
	fmt.Println("The status trie can be built")

	return bc, nil
}

// 本地化存储，修改内存、存储至硬盘。返回要迁出账户的余额
func (bc *BlockChain) AddBlock(block *core.Block) map[string]*big.Int {
	iobegin := time.Now().UnixMilli()
	bc.Storage.AddBlock(block)
	var outbalance map[string]*big.Int

	if !bc.ChainConfig.Stop_When_Migrating {
		// 将迁入账户映射到本分片
		for _, v := range block.TXmig2s {
			account.Account2ShardLock.Lock()
			account.Account2Shard[v.Address] = params.ShardTable[params.Config.ShardID]
			account.AccountInOwnShard[v.Address] = true
			account.Account2ShardLock.Unlock()
		}
	}

	updatetree := time.Now().UnixMilli()
	stateroothash, outbalance := bc.getUpdatedTreeOfState(1, block.Header.Number, block.Transactions, block.TXmig1s, block.TXmig2s, block.Anns, block.NSs)
	if !bytes.Equal(block.Header.StateRoot, stateroothash){
		log.Panicf("二者不等，Ins长度为%v\n", len(block.TXmig2s))
	}
	fmt.Printf("更新树花时间为: %v\n", time.Now().UnixMilli()-updatetree)

	if !bc.ChainConfig.Stop_When_Migrating && params.Config.NodeID != "N0" {
		// 将删除账户从本分片映射删去，并设为对应分片
		for _, v := range block.Anns {
			account.Account2ShardLock.Lock()
			account.Account2Shard[v.Address] = v.ToshardID
			delete(account.AccountInOwnShard, v.Address)
			account.Account2ShardLock.Unlock()

			if bc.ChainConfig.Lock_Acc_When_Migrating {
				account.Lock_Acc_Lock.Lock()
				delete(account.Lock_Acc, v.Address)
				account.Lock_Acc_Lock.Unlock()
			}
		}
	}

	// 若要锁账户，就把账户锁住
	if params.Config.Lock_Acc_When_Migrating && params.Config.NodeID != "N0" {
		account.Lock_Acc_Lock.Lock()
		for _, v := range block.TXmig1s {
			account.Lock_Acc[v.Address] = true
			// if params.Config.NodeID == "N0" {
			// 	bc.Tx_pool.Locking_TX_Pools[k] = make([]*core.Transaction, 0)
			// }
		}
		account.Lock_Acc_Lock.Unlock()
	}

	if params.Config.Fail && params.Config.Fail_Time+1 == block.Header.Number && params.Config.Lock_Acc_When_Migrating {
		account.Lock_Acc_Lock.Lock()
		account.Lock_Acc["489338d5e8d42e8c923d1f47361d979503d4ad68"] = false
		account.Lock_Acc_Lock.Unlock()
	}

	bc.CurrentBlock = block

	fmt.Printf("写数据库IO花时间为: %v\n", time.Now().UnixMilli()-iobegin)

	// relay
	// if bc.ChainConfig.NodeID == "N0" {
	// 	bc.genRelayTxs(block)
	// }
	return outbalance
}

// func (bc *BlockChain) genRelayTxs(block *core.Block) {
// 	for _, tx := range block.Transactions {
// 		shardID := account.Addr2Shard(hex.EncodeToString(tx.Recipient))
// 		if shardID != params.ShardTable[bc.ChainConfig.ShardID] {
// 			bc.Tx_pool.AddRelayTx(tx, params.ShardTableInt2Str[shardID])
// 		}
// 	}
// }

func (bc *BlockChain) GenerateBlock(id int) *core.Block {
	quota := params.Config.MaxMigSize
	mig1s := []*core.TXmig1{}
	mig2s := []*core.TXmig2{}
	anns := []*core.TXann{}
	// deleted := []*core.Delete{}
	nss := []*core.TXns{}
	if !params.Config.Stop_When_Migrating {
		//得到新的映射（只有要改变的

		if params.Config.Algorithm || params.Config.Pressure {
			mig1s = bc.TXmig1_pool.FetchTXmig1s2Pack2()
		} else {
			mig1s, quota = bc.TXmig1_pool.FetchTXmig1s2Pack()
		}

		// bc.Tx_pool.Lock.Lock()
		// new_migration := make(map[string]int)
		// tmp, _ := json.Marshal(&bc.Tx_pool.Migration_Pool)
		// _ = json.Unmarshal(tmp, &new_migration)
		// bc.Tx_pool.Migration_Pool = make(map[string]int)
		// bc.Tx_pool.Lock.Unlock()

		// account.Account2ShardLock.Lock()
		// //找到要迁移出去的
		// for addr, shard := range new_migration {
		// 	if account.AccountInOwnShard[addr] && shard != params.ShardTable[params.Config.ShardID] {
		// 		out[addr] = shard
		// 	}
		// }
		// account.Account2ShardLock.Unlock()

		if !params.Config.Bu_Tong_Shi_Jian {
			if params.Config.Algorithm || params.Config.Pressure {
				mig2s = bc.TXmig2_pool.FetchTXmig2s2Pack2()
			} else {
				//将全部迁入请求包进来
				mig2s, quota = bc.TXmig2_pool.FetchTXmig2s2Pack(quota)
			}
		} else if id == params.Config.Bu_Tong_Shi_Jian_Jian_Ge {
			//将全部迁入请求包进来
			mig2s, _ = bc.TXmig2_pool.FetchTXmig2s2Pack(10000000)
		}

		// //将全部该删除账户包进来
		// if params.Config.Algorithm || params.Config.Pressure {
		// 	deleted = bc.Delete_pool.FetchDels2Pack2()
		// } else {
		// 	deleted = bc.Delete_pool.FetchDels2Pack()
		// }

		//将全部announce包进来
		if params.Config.Algorithm || params.Config.Pressure {
			anns = bc.TXann_pool.FetchTXanns2Pack2()
		} else {
			anns, quota = bc.TXann_pool.FetchTXanns2Pack(quota)
		}

		//将全部加（减）钱请求包进来
		if params.Config.Algorithm || params.Config.Pressure {
			nss = bc.TXns_pool.FetchTXnss2Pack2()
		} else {
			nss, quota = bc.TXns_pool.FetchTXnss2Pack(quota)
		}
	}

	txs := []*core.Transaction{}
	queueLen := 0
	if params.Config.Cross_Chain && id != 1 {
		if params.Config.ShardID == "S0" && (params.Config.Fail || (!params.Config.Fail && !params.Config.Lock_Acc_When_Migrating)) {
			txs, queueLen = bc.Tx_pool.FetchTxs2Pack(params.Config.MaxBlockSize, bc.CurrentBlock.Header.Number+1)
		} else if params.Config.Lock_Acc_When_Migrating && params.Config.ShardID == "S1" && !params.Config.Fail {
			account.Account2ShardLock.Lock()
			if account.AccountInOwnShard["489338d5e8d42e8c923d1f47361d979503d4ad68"] {
				account.Account2ShardLock.Unlock()
				txs, queueLen = bc.Tx_pool.FetchTxs2Pack(params.Config.MaxBlockSize, bc.CurrentBlock.Header.Number+1)
			} else {
				account.Account2ShardLock.Unlock()
			}
		}
	} else if (!params.Config.Bu_Tong_Bi_Li && !params.Config.Bu_Tong_Shi_Jian && !params.Config.Fail && !params.Config.Cross_Chain) || id != 1 {
		if params.Config.Algorithm || params.Config.Pressure {
			txs, queueLen = bc.Tx_pool.FetchTxs2Pack(params.Config.MaxBlockSize - len(mig1s) - len(mig2s) - len(anns) - len(nss), bc.CurrentBlock.Header.Number + 1)
		} else {
			txs, queueLen = bc.Tx_pool.FetchTxs2Pack(params.Config.MaxBlockSize - len(mig1s) - len(mig2s) - len(anns) - len(nss), bc.CurrentBlock.Header.Number + 1)
			// txs = bc.Tx_pool.FetchTxs2Pack(params.Config.MaxBlockSize - params.Config.MaxMigSize + quota)
		}
	}

	blockHeader := &core.BlockHeader{
		ParentHash: bc.CurrentBlock.Hash,
		Number:     bc.CurrentBlock.Header.Number + 1,
		Time:       uint64(time.Now().Unix()),
	}

	block := core.NewBlock(blockHeader, txs, mig1s, mig2s, anns, nss)

	//普通交易树
	block.Header.TxHash = GetTxTreeRoot(txs)

	//迁移交易树
	block.Header.MigHash = GetMigTreeRoot(mig1s, mig2s, anns, nss)

	//0 代表不更新到磁盘
	block.Header.StateRoot, _ = bc.getUpdatedTreeOfState(0, blockHeader.Number, txs, mig1s, mig2s, anns, nss)

	block.Hash = block.GetHash()

	if !params.Config.Stop_When_Migrating {
		if !params.Config.Lock_Acc_When_Migrating {
			account.Outing_Acc_Before_Announce_Lock.Lock()
			for _,lockedtxs := range bc.Tx_pool.Outing_Before_Announce_TX_Pools {
				queueLen += len(lockedtxs)
			}
			account.Outing_Acc_Before_Announce_Lock.Unlock()
		}else {
			account.Lock_Acc_Lock.Lock()
			for _,lockedtxs := range bc.Tx_pool.Locking_TX_Pools {
				queueLen += len(lockedtxs)
			}
			account.Lock_Acc_Lock.Unlock()
		}

	}
	s := fmt.Sprintf("%v %v", blockHeader.Number, queueLen)
	queuelenlog.Write(strings.Split(s, " "))
	queuelenlog.Flush()

	return block
}

// 输出更新后的状态树 以及 要迁出的账户的余额
func (bc *BlockChain) getUpdatedTreeOfState(commit int, height int, txs []*core.Transaction, mig1s []*core.TXmig1, mig2s []*core.TXmig2, anns []*core.TXann, nss []*core.TXns) ([]byte, map[string]*big.Int) {
	// build trie from the triedb (in disk)
	start_execute := time.Now().UnixMicro()
	st, err := trie.New(trie.TrieID(common.BytesToHash(bc.CurrentBlock.Header.StateRoot)), bc.Triedb)
	if err != nil {
		log.Panic(err)
	}

	mig2start := time.Now().UnixMicro()
	if !bc.ChainConfig.Stop_When_Migrating {
		// 将迁入账户加到状态树
		for _, v := range mig2s {
			hex_address, _ := hex.DecodeString(v.Address)
			s_state_enc := st.Get(hex_address)
			if s_state_enc == nil {
				log.Panic()
			}
			account_state := account.DecodeAccountState(s_state_enc)
			account_state.Balance.Set(v.Value)
			account_state.Migrate = -1
			account_state.Location = params.ShardTable[bc.ChainConfig.ShardID]
			st.Update(hex_address, account_state.Encode())
		}
	}
	mig2time := time.Now().UnixMicro() - mig2start

	// stateTree := bc.preExecute(txs)
	outbalance := make(map[string]*big.Int)

	txstart := time.Now().UnixMicro()
	for _, tx := range txs {
		// 确保发送地址属于此分片，即此交易不是其它分片发来的relay交易
		// if account.Addr2Shard(hex.EncodeToString(tx.Sender)) == params.ShardTable[bc.ChainConfig.ShardID] {
		if !tx.IsRelay && !tx.Relay_Lock {
			s_state_enc := st.Get(tx.Sender)
			if s_state_enc == nil {
				fmt.Printf("sender属于该分片吗：%v\n", account.AccountInOwnShard[hex.EncodeToString(tx.Sender)])
				fmt.Printf("sender地址为：%v\n", hex.EncodeToString(tx.Sender))
				fmt.Printf("sender属于分片：%v\n", account.Account2Shard[hex.EncodeToString(tx.Sender)])
				fmt.Printf("rec地址为：%v\n", hex.EncodeToString(tx.Recipient))
				fmt.Printf("rec属于分片：%v\n", account.Account2Shard[hex.EncodeToString(tx.Recipient)])
				log.Panic()
			}
			account_state := account.DecodeAccountState(s_state_enc)
			account_state.Balance.Sub(account_state.Balance, tx.Value)
			st.Update(tx.Sender, account_state.Encode())
		}


		r_state_enc := st.Get(tx.Recipient)
		if r_state_enc == nil {
			fmt.Printf("rec属于该分片吗：%v\n", account.AccountInOwnShard[hex.EncodeToString(tx.Recipient)])
			fmt.Printf("rec地址为：%v\n", hex.EncodeToString(tx.Recipient))
			fmt.Printf("rec属于分片：%v\n", account.Account2Shard[hex.EncodeToString(tx.Recipient)])
			fmt.Printf("txid：%v\n", tx.Id)
			log.Panic()
		}
		account_state := account.DecodeAccountState(r_state_enc)
		// 接收地址不在此分片，不对该状态进行修改
		if account_state.Location != params.ShardTable[bc.ChainConfig.ShardID] {
			continue
			// 接收者为锁定账户，不对该状态进行修改
		} else if account_state.Migrate != -1 && params.Config.Lock_Acc_When_Migrating{
			continue
		}

		// 接收者为锁定账户，不对该状态进行修改
		// if params.Config.RelayLock && tx.Rec_Suppose_on_chain == height {
		// 	continue
		// }

		// 接收者为锁定账户，不对该状态进行修改
		// if params.Config.Lock_Acc_When_Migrating {
		// 	account.Lock_Acc_Lock.Lock()
		// 	if account.Lock_Acc[hex.EncodeToString(tx.Recipient)] {
		// 		account.Lock_Acc_Lock.Unlock()
		// 		continue
		// 	}
		// 	account.Lock_Acc_Lock.Unlock()
		// }

		
		account_state.Balance.Add(account_state.Balance, tx.Value)
		if commit == 1 && account_state.Migrate != -1 {
			tx.HalfLock = true
		}
		st.Update(tx.Recipient, account_state.Encode())

	}
	txtime := time.Now().UnixMicro() - txstart

	mig1time := int64(0)
	anntime := int64(0)
	nstime := int64(0)
	if !bc.ChainConfig.Stop_When_Migrating {
		// 状态树中记录要迁出的账户的去向
		mig1start := time.Now().UnixMicro()
		for _, v := range mig1s {
			hex_address, _ := hex.DecodeString(v.Address)
			encoded := st.Get(hex_address)
			if encoded == nil {
				log.Panic()
			}
			account_state := account.DecodeAccountState(encoded)
			account_state.Migrate = v.ToshardID
			outbalance[v.Address] = new(big.Int).Set(account_state.Balance)
			st.Update(hex_address, account_state.Encode())
		}
		mig1time = time.Now().UnixMicro() - mig1start

		// 状态树中记录announce的账户
		annstart := time.Now().UnixMicro()
		for _, v := range anns {
			hex_address, _ := hex.DecodeString(v.Address)
			s_state_enc := st.Get(hex_address)
			if s_state_enc == nil {
				log.Panic()
			}
			account_state := account.DecodeAccountState(s_state_enc)
			account_state.Migrate = -1
			account_state.Location = v.ToshardID
			st.Update(hex_address, account_state.Encode())
			// encoded := st.Get(hex_address)
			// account_state := &account.AccountState{
			// 	Migrate:  -1,
			// 	Location: v.ToshardID,
			// }
			// if encoded != nil {
			// 	account_state = account.DecodeAccountState(encoded)
			// 	account_state.Migrate = -1
			// 	account_state.Location = v.ToshardID
			// }

			// st.Update(hex_address, account_state.Encode())
		}
		anntime = time.Now().UnixMicro() - annstart

		// // 将迁入账户加到状态树
		// for _, v := range in1s {
		// 	hex_address, _ := hex.DecodeString(v.Address)

		// 	account_state := &account.AccountState{
		// 		Balance: v.Value,
		// 		Migrate: -1,
		// 	}
		// 	stateTree.Put(hex_address, account_state.Encode())
		// }

		// // 从状态树中删除指定账户
		// for _, v := range deletes {
		// 	hex_address, _ := hex.DecodeString(v.Address)
		// 	stateTree.Delete(hex_address)
		// }

		// 将加减钱改到状态树
		nsstart := time.Now().UnixMicro()
		for _, v := range nss {
			hex_address, _ := hex.DecodeString(v.Address)

			encoded := st.Get(hex_address)
			if encoded == nil {
				log.Panic()
			}
			account_state := account.DecodeAccountState(encoded)
			account_state.Balance.Add(account_state.Balance, v.Change)
			st.Update(hex_address, account_state.Encode())
		}
		nstime = time.Now().UnixMicro() - nsstart
	}

	st_hash_bytes := st.Hash().Bytes()
	//只是生成区块
	if commit == 0 {
		return st_hash_bytes, outbalance
	}

	//区块上链，要写到磁盘
	// commit the memory trie to the database in the disk
	rt, ns := st.Commit(false)
	//空块
	if ns == nil {
		blocktime := time.Now().UnixMicro() - start_execute
		if commit == 1 && params.Config.NodeID == "N0" {
			s := fmt.Sprintf("%v %v %v %v %v %v %v", height, blocktime, txtime, mig1time, mig2time, anntime, nstime)
			blocktimelog.Write(strings.Split(s, " "))
			blocktimelog.Flush()
		}
		return st_hash_bytes, outbalance
	}
	err = bc.Triedb.Update(trie.NewWithNodeSet(ns))
	if err != nil {
		log.Panic()
	}
	err = bc.Triedb.Commit(rt, false)
	if err != nil {
		log.Panic(err)
	}
	blocktime := time.Now().UnixMicro() - start_execute
	if commit == 1 && params.Config.NodeID == "N0" {
		s := fmt.Sprintf("%v %v %v %v %v %v %v", height, blocktime, txtime, mig1time, mig2time, anntime, nstime)
		blocktimelog.Write(strings.Split(s, " "))
		blocktimelog.Flush()
	}
	return rt.Bytes(), outbalance
}

func (bc *BlockChain) NewGenesisBlock() *core.Block {
	blockHeader := &core.BlockHeader{
		Number: 0,
		Time:   uint64(time.Date(2022, 05, 28, 17, 11, 0, 0, time.Local).Unix()),
	}

	txs := make([]*core.Transaction, 0)
	mig1s := make([]*core.TXmig1, 0)
	mig2s := make([]*core.TXmig2, 0)
	anns := make([]*core.TXann, 0)
	nss := make([]*core.TXns, 0)
	block := core.NewBlock(blockHeader, txs, mig1s, mig2s, anns, nss)

	// build a new trie database by db
	triedb := trie.NewDatabaseWithConfig(bc.db, &trie.Config{
		Cache:     0,
		Preimages: true,
	})
	bc.Triedb = triedb
	statusTrie := trie.NewEmpty(triedb)
	block.Header.StateRoot = bc.genesisStateTree(statusTrie.Hash().Bytes())
	block.Header.TxHash = GetTxTreeRoot(txs)
	block.Hash = block.GetHash()

	return block
}

func (bc *BlockChain) AddGenesisBlock(block *core.Block) {
	bc.Storage.AddBlock(block)

	// 重新从数据库中获取最新内容
	newestBlockHash, err := bc.Storage.GetNewestBlockHash()
	if err != nil {
		log.Panic()
	}
	curBlock, err := bc.Storage.GetBlock(newestBlockHash)
	if err != nil {
		log.Panic()
	}
	bc.CurrentBlock = curBlock
}

// 创世区块中初始化几个账户
// func genesisStateTree() *trie.Trie {
// 	trie := trie.NewTrie()
// 	for i := 0; i < len(params.Init_addrs); i++ {
// 		address := params.Init_addrs[i]
// 		if utils.Addr2Shard(address) != params.ShardTable[params.Config.ShardID] {
// 			continue
// 		}
// 		value, _ := strconv.ParseFloat(params.Init_balance, 64)
// 		accountState := &account.AccountState{
// 			Balance: value,
// 			Migrate: -1,
// 		}
// 		hex_address, _ := hex.DecodeString(address)
// 		trie.Put(hex_address, accountState.Encode())
// 	}
// 	return trie
// }

// 创世区块中初始化几个账户
func (bc *BlockChain) genesisStateTree(stateroot []byte) []byte {
	// build trie from the triedb (in disk)
	st, err := trie.New(trie.TrieID(common.BytesToHash(stateroot)), bc.Triedb)
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < len(params.Init_addrs); i++ {
		address := params.Init_addrs[i]
		// if utils.Addr2Shard(address) != params.ShardTable[params.Config.ShardID] {
		// 	continue
		// }
		value := new(big.Int)
		value, ok := value.SetString(params.Init_balance, 10)
		if !ok {
			log.Panic()
		}
		accountState := &account.AccountState{
			Balance:  value,
			Migrate:  -1,
			Location: utils.Addr2Shard(address),
		}
		hex_address, _ := hex.DecodeString(address)
		st.Update(hex_address, accountState.Encode())
	}
	// commit the memory trie to the database in the disk
	rt, ns := st.Commit(false)
	err = bc.Triedb.Update(trie.NewWithNodeSet(ns))
	if err != nil {
		log.Panic()
	}
	err = bc.Triedb.Commit(rt, false)
	if err != nil {
		log.Panic(err)
	}
	return rt.Bytes()
}

// Get the transaction root, this root can be used to check the transactions
func GetTxTreeRoot(txs []*core.Transaction) []byte {
	// use a memory trie database to do this, instead of disk database
	triedb := trie.NewDatabase(rawdb.NewMemoryDatabase())
	transactionTree := trie.NewEmpty(triedb)
	for _, tx := range txs {
		transactionTree.Update(tx.TxHash, tx.Encode())
	}
	return transactionTree.Hash().Bytes()
}

// Get the transaction root, this root can be used to check the transactions
func GetMigTreeRoot(mig1s []*core.TXmig1, mig2s []*core.TXmig2, anns []*core.TXann, nss []*core.TXns) []byte {
	// use a memory trie database to do this, instead of disk database
	triedb := trie.NewDatabase(rawdb.NewMemoryDatabase())
	transactionTree := trie.NewEmpty(triedb)
	for _, tx := range mig1s {
		transactionTree.Update(tx.Hash(), tx.Encode())
	}
	for _, tx := range mig2s {
		transactionTree.Update(tx.Hash(), tx.Encode())
	}
	for _, tx := range anns {
		transactionTree.Update(tx.Hash(), tx.Encode())
	}
	for _, tx := range nss {
		transactionTree.Update(tx.Hash(), tx.Encode())
	}
	return transactionTree.Hash().Bytes()
}

func (bc *BlockChain) IsBlockValid(block *core.Block) bool {
	// todo

	return true
}
