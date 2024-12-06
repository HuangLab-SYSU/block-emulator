package chain

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"blockEmulator/core"
	"blockEmulator/new_trie"
	"blockEmulator/params"
	"blockEmulator/storage"
	"blockEmulator/utils"
)

type BlockChain struct {
	ChainConfig *params.ChainConfig // Chain configuration

	CurrentBlock *core.Block // Current head of the block chain

	Storage *storage.Storage

	StatusTrie *new_trie.N_Trie

	Tx_pool *core.Tx_pool
}

func NewBlockChain(chainConfig *params.ChainConfig) (*BlockChain, error) {
	fmt.Printf("%v\n", chainConfig)

	bc := &BlockChain{
		ChainConfig: chainConfig,
		Storage:     storage.NewStorage(chainConfig),
		Tx_pool:     core.NewTxPool(),
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

	block, err := bc.Storage.GetBlock(blockHash)
	if err != nil {
		log.Panic()
	}
	bc.CurrentBlock = block

	stateTree, err := bc.Storage.GetStatusTree()
	if err != nil {
		log.Panic()
	}
	bc.StatusTrie = stateTree

	return bc, nil
}
func (bc *BlockChain) NewBlockChainAfterReconfig(chainConfig *params.ChainConfig) {

	// bc.Storage = storage.NewStorage(chainConfig)
	fmt.Printf("%v\n", chainConfig)

	bc = &BlockChain{
		ChainConfig: chainConfig,
		Storage:     storage.NewStorage(chainConfig),
		Tx_pool:     core.NewTxPool(),
	}

	blockHash, err := bc.Storage.GetNewestBlockHash()
	if err != nil { // 应该返回nil，即不存在newestHash
		if err.Error() == "newestBlockHash is not found" {
			genesisBlock := bc.NewGenesisBlock()
			bc.AddGenesisBlock(genesisBlock)

		}
		fmt.Println("=====1=======")
		// log.Panic()
	}

	block, err := bc.Storage.GetBlock(blockHash)
	if err != nil {
		fmt.Println("=====2=======")
		// log.Panic()
	}
	bc.CurrentBlock = block

	stateTree, err := bc.Storage.GetStatusTree()
	if err != nil {
		fmt.Println("=====3=======")
		// log.Panic()
	}
	bc.StatusTrie = stateTree
	fmt.Printf("======更新完毕=======\n")

}

func (bc *BlockChain) AddBlock(block *core.Block, consensusFlag bool, epochID int) int {
	num := 0
	bc.Storage.AddBlock(block)
	bc.StatusTrie, num = bc.getUpdatedTreeOfState(block.Transactions, epochID) // 更新状态树trie（内存）
	bc.Storage.UpdateStateTree(bc.StatusTrie)                                  // 根据trie更新statetreebucket（外存）

	// // 重新从数据库中获取最新内容
	// newestBlockHash, err := bc.Storage.GetNewestBlockHash()
	// if err != nil {
	// 	log.Panic()
	// }
	// curBlock, err := bc.Storage.GetBlock(newestBlockHash)
	// if err != nil {
	// 	log.Panic()
	// }
	bc.CurrentBlock = block
	// stateTree, err := bc.Storage.GetStatusTree()
	// if err != nil {
	// 	log.Panic()
	// }
	// bc.StatusTrie = stateTree

	// 只有在共识周期上链时增加relay txs
	if consensusFlag {
		if bc.ChainConfig.NodeID == "N0" {
			bc.genRelayTxs(block)
		}
	}
	return num
}

func (bc *BlockChain) genRelayTxs(block *core.Block) {
	for _, tx := range block.Transactions {
		shardID := utils.Addr2Shard(hex.EncodeToString(tx.Recipient))
		if shardID != params.ShardTable[bc.ChainConfig.ShardID] {
			bc.Tx_pool.AddRelayTx(tx, params.ShardTableInt2Str[shardID])
		}
	}

}

func (bc *BlockChain) GenerateBlock() (*core.Block, int, int, []string, []string) {
	fmt.Printf("len(pool.Queue): %d\n", len(bc.Tx_pool.Queue))

	txs := bc.Tx_pool.FetchTxs2Pack()
	blockHeader := &core.BlockHeader{
		ParentHash: bc.CurrentBlock.Hash,
		Number:     bc.CurrentBlock.Header.Number + 1,
		Time:       uint64(time.Now().Unix()),
	}

	block := core.NewBlock(blockHeader, txs)

	mpt_tx := bc.buildTreeOfTxs(txs)
	block.Header.TxHash = mpt_tx.Hash()

	mpt_state, _ := bc.getUpdatedTreeOfState(txs, 0) // epochID无用，用0填充
	block.Header.StateRoot = mpt_state.Hash()

	block.Hash = block.GetHash()

	// add 找到需要请求的地址
	var accounts_send []string
	var accounts_all []string
	var num_accounts_visit int
	var num_active_accounts_visit int

	stateTree := &new_trie.N_Trie{}
	new_trie.DeepCopy(stateTree, bc.StatusTrie)

	for _, tx := range txs {
		sender := tx.Sender
		if utils.Addr2Shard(hex.EncodeToString(sender)) == params.ShardTable[bc.ChainConfig.ShardID] { // 发送地址在此分片，此交易不是其它分片发送过来的relay交易
			num_accounts_visit += 1
			if _, ok := stateTree.Get(sender); !ok { // 若原状态树中不存在发送账户
				flag := false
				for _, account := range accounts_send {
					if account == string(sender) {
						flag = true
						break
					}
				}
				if !flag {
					accounts_send = append(accounts_send, string(sender))
				}

			} else {
				num_active_accounts_visit += 1
			}
			flag := false
			for _, account := range accounts_all {
				if account == string(sender) {
					flag = true
					break
				}
			}
			if !flag {
				accounts_all = append(accounts_all, string(sender))
			}
		}

		receiever := tx.Recipient
		num_accounts_visit += 1
		if _, ok := stateTree.Get(receiever); !ok { // 若原状态树中不存在发送账户
			flag := false
			for _, account := range accounts_send {
				if account == string(receiever) {
					flag = true
					break
				}
			}
			if !flag {
				accounts_send = append(accounts_send, string(receiever))
			}

			// accounts_send = append(accounts_send, string(receiever))
		} else {
			num_active_accounts_visit += 1
		}

		flag := false
		for _, account := range accounts_all {
			if account == string(receiever) {
				flag = true
				break
			}
		}
		if !flag {
			accounts_all = append(accounts_all, string(receiever))
		}
	}

	return block, num_accounts_visit, num_active_accounts_visit, accounts_send, accounts_all
}

func (bc *BlockChain) buildTreeOfTxs(txs []*core.Transaction) *new_trie.N_Trie {
	trie := new_trie.NewTrie()
	for _, tx := range txs {
		trie.Put(tx.TxHash[:], tx.Encode(), 0)
	}
	return trie
}

func (bc *BlockChain) getUpdatedTreeOfState(txs []*core.Transaction, epochID int) (*new_trie.N_Trie, int) {
	stateTree, num := bc.preExecute(txs, bc, epochID) // 对状态树trie中不存在的账户，创建空节点
	// stateTree, err := bc.Storage.GetStatusTree()
	// if err != nil {
	// 	if err.Error() == "stateTree is not found" {
	// 		stateTree = genesisStateTree()
	// 	} else {
	// 		log.Panic()
	// 	}
	// }

	for _, tx := range txs {
		// 确保发送地址属于此分片，即此交易不是其它分片发来的relay交易
		if utils.Addr2Shard(hex.EncodeToString(tx.Sender)) == params.ShardTable[bc.ChainConfig.ShardID] {
			decoded, success := stateTree.Get(tx.Sender)
			if !success {
				log.Panic()
			}
			account_state := core.DecodeAccountState(decoded)
			account_state.Deduct(tx.Value)
			stateTree.Put(tx.Sender, account_state.Encode(), epochID) // 更新交易树
		}

		// 接收地址不在此分片，不对该状态进行修改
		if utils.Addr2Shard(hex.EncodeToString(tx.Recipient)) != params.ShardTable[bc.ChainConfig.ShardID] {
			continue
		}
		decoded, success := stateTree.Get(tx.Recipient)
		if !success {
			log.Panic()
		}
		account_state := core.DecodeAccountState(decoded)
		account_state.Deposit(tx.Value)
		stateTree.Put(tx.Recipient, account_state.Encode(), epochID)

	}
	return stateTree, num
}

func (bc *BlockChain) NewGenesisBlock() *core.Block {
	blockHeader := &core.BlockHeader{
		Number: 0,
		Time:   uint64(time.Date(2022, 11, 28, 17, 11, 0, 0, time.Local).Unix()),
	}

	txs := make([]*core.Transaction, 0)
	block := core.NewBlock(blockHeader, txs)

	stateTree := genesisStateTree()
	bc.Storage.UpdateStateTree(stateTree)
	block.Header.StateRoot = stateTree.Hash()

	mpt_tx := bc.buildTreeOfTxs(txs)
	block.Header.TxHash = mpt_tx.Hash()

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
	stateTree, err := bc.Storage.GetStatusTree()
	if err != nil {
		log.Panic()
	}
	bc.StatusTrie = stateTree
}

// 创世区块中初始化几个账户
func genesisStateTree() *new_trie.N_Trie {
	trie := new_trie.NewTrie()
	for i := 0; i < len(params.Init_addrs); i++ {
		address := params.Init_addrs[i]
		value := new(big.Int)
		value.SetString(params.Init_balance, 10)
		accountState := &core.AccountState{
			Balance: value,
		}
		hex_address, _ := hex.DecodeString(address)
		trie.Put(hex_address, accountState.Encode(), 0)
	}
	return trie
}

// func GenerateTxs(bc *BlockChain) {
// 	tx_cnt := 2
// 	txs := make([]*core.Transaction, tx_cnt)
// 	for i := 0; i < tx_cnt; i++ {
// 		sender := core.GenerateAddress()
// 		receiver := core.GenerateAddress()
// 		txs[i] = &core.Transaction{
// 			Sender:    sender,
// 			Recipient: receiver,
// 			Value:     big.NewInt(int64(i + 1)),
// 		}
// 	}
// 	bc.Tx_pool.AddTxs(txs)
// 	bc.preExecute(txs, bc)
// }

// 不改变原本状态树
// 返回状态树与新增节点个数
func (bc *BlockChain) preExecute(txs []*core.Transaction, chain *BlockChain, epochID int) (*new_trie.N_Trie, int) {
	// stateTree, err := chain.Storage.GetStatusTree()
	// if err != nil {
	// 	log.Panic()
	// }
	stateTree := &new_trie.N_Trie{}
	new_trie.DeepCopy(stateTree, bc.StatusTrie)

	num := 0

	for _, tx := range txs {
		sender := tx.Sender
		if utils.Addr2Shard(hex.EncodeToString(sender)) == params.ShardTable[bc.ChainConfig.ShardID] { // 发送地址在此分片，此交易不是其它分片发送过来的relay交易
			if _, ok := stateTree.Get(sender); !ok { // 若原状态树中不存在发送账户
				num += 1
				value := new(big.Int)
				value.SetString(params.Init_balance, 10)
				accountState := &core.AccountState{
					Balance: value,
				}
				stateTree.Put(sender, accountState.Encode(), epochID)
			}
		}

		receiver := tx.Recipient
		if utils.Addr2Shard(hex.EncodeToString(receiver)) != params.ShardTable[bc.ChainConfig.ShardID] { // 接收地址不在此分片
			continue
		}
		if _, ok := stateTree.Get(receiver); !ok { // 若原状态树中不存在接收账户
			num += 1
			value := new(big.Int)
			value.SetString(params.Init_balance, 10)
			accountState := &core.AccountState{
				Balance: value,
			}
			stateTree.Put(receiver, accountState.Encode(), epochID)
		}
	}
	// chain.Storage.UpdateStateTree(stateTree)
	return stateTree, num

}

func (bc *BlockChain) IsBlockValid(block *core.Block) bool {
	// todo

	return true
}
