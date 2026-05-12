package committee

import (
	"blockEmulator/broker"
	"blockEmulator/core"
	"blockEmulator/message"
	"blockEmulator/networks"
	"blockEmulator/params"
	"blockEmulator/supervisor/Broker2Earn"
	"blockEmulator/supervisor/signal"
	"blockEmulator/supervisor/supervisor_log"
	"blockEmulator/utils"
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// CLPA committee operations
type BrokerCommitteeMod_b2e struct {
	csvPath      string
	dataTotalNum int
	nowDataNum   int
	dataTxNums   int
	batchDataNum int

	//broker related  attributes avatar
	broker               *broker.Broker
	brokerConfirm1Pool   map[string]*message.Mag1Confirm
	brokerConfirm2Pool   map[string]*message.Mag2Confirm
	restBrokerRawMegPool []*message.BrokerRawMeg
	brokerTxPool         []*core.Transaction
	brokerModuleLock     sync.Mutex
	brokerBalanceLock    sync.Mutex

	// logger module
	sl *supervisor_log.SupervisorLog

	// control components
	Ss          *signal.StopSignal // to control the stop message sending
	IpNodeTable map[uint64]map[uint64]string

	// log balance
	Result_lockBalance   map[string][]string
	Result_brokerBalance map[string][]string
	Result_Profit        map[string][]string

	// batch injection control: block next batch until current batch is fully confirmed on-chain
	b2eLock             sync.Mutex
	collectTransactions float64

	// epoch stats recording
	epochBroker1TxNum  int
	epochBroker2TxNum  int
	epochInnerTxNum    int
	epochBlockCount    int
	epochAllocatedTx   int
	epochUnallocatedTx int
	totalB2EIterations int
}

func NewBrokerCommitteeMod_b2e(Ip_nodeTable map[uint64]map[uint64]string, Ss *signal.StopSignal, sl *supervisor_log.SupervisorLog, csvFilePath string, dataNum, batchNum int) *BrokerCommitteeMod_b2e {

	broker := new(broker.Broker)
	broker.NewBroker(nil)
	result_lockBalance := make(map[string][]string)
	result_brokerBalance := make(map[string][]string)
	result_Profit := make(map[string][]string)
	block_txs := make(map[uint64][]string)

	for _, brokeraddress := range broker.BrokerAddress {
		result_lockBalance[brokeraddress] = make([]string, 0)
		result_brokerBalance[brokeraddress] = make([]string, 0)
		result_Profit[brokeraddress] = make([]string, 0)

		a := ""
		b := ""
		title := ""
		for i := 0; i < params.ShardNum; i++ {
			title += "shard" + strconv.Itoa(i) + ","
			a += params.Init_broker_Balance.String() + ","
			b += "0,"
		}
		result_lockBalance[brokeraddress] = append(result_lockBalance[brokeraddress], title)
		result_brokerBalance[brokeraddress] = append(result_brokerBalance[brokeraddress], title)
		result_Profit[brokeraddress] = append(result_Profit[brokeraddress], title)

		result_lockBalance[brokeraddress] = append(result_lockBalance[brokeraddress], b)
		result_brokerBalance[brokeraddress] = append(result_brokerBalance[brokeraddress], a)
		result_Profit[brokeraddress] = append(result_Profit[brokeraddress], b)
	}
	for i := 0; i < params.ShardNum; i++ {
		block_txs[uint64(i)] = make([]string, 0)
		block_txs[uint64(i)] = append(block_txs[uint64(i)], "txExcuted, broker1Txs, broker2Txs, allocatedTxs")
	}

	return &BrokerCommitteeMod_b2e{
		csvPath:              csvFilePath,
		dataTotalNum:         dataNum,
		batchDataNum:         batchNum,
		nowDataNum:           0,
		dataTxNums:           0,
		brokerConfirm1Pool:   make(map[string]*message.Mag1Confirm),
		brokerConfirm2Pool:   make(map[string]*message.Mag2Confirm),
		restBrokerRawMegPool: make([]*message.BrokerRawMeg, 0),
		brokerTxPool:         make([]*core.Transaction, 0),
		broker:               broker,
		IpNodeTable:          Ip_nodeTable,
		Ss:                   Ss,
		sl:                   sl,
		Result_lockBalance:   result_lockBalance,
		Result_brokerBalance: result_brokerBalance,
		Result_Profit:        result_Profit,
	}

}

func (bcm *BrokerCommitteeMod_b2e) HandleOtherMessage([]byte) {}

func (bcm *BrokerCommitteeMod_b2e) fetchModifiedMap(key string) uint64 {
	return uint64(utils.Addr2Shard(key))
}

func (bcm *BrokerCommitteeMod_b2e) txSending(txlist []*core.Transaction) {
	// the txs will be sent
	sendToShard := make(map[uint64][]*core.Transaction)

	for idx := 0; idx <= len(txlist); idx++ {
		if idx > 0 && (idx%params.InjectSpeed == 0 || idx == len(txlist)) {
			// send to shard
			for sid := uint64(0); sid < uint64(params.ShardNum); sid++ {
				it := message.InjectTxs{
					Txs:       sendToShard[sid],
					ToShardID: sid,
				}
				itByte, err := json.Marshal(it)
				if err != nil {
					log.Panic(err)
				}
				send_msg := message.MergeMessage(message.CInject, itByte)
				go networks.TcpDial(send_msg, bcm.IpNodeTable[sid][0])
			}
			sendToShard = make(map[uint64][]*core.Transaction)
			time.Sleep(time.Second)
		}
		if idx == len(txlist) {
			break
		}
		tx := txlist[idx]
		sendersid := bcm.fetchModifiedMap(tx.Sender)

		if bcm.broker.IsBroker(tx.Sender) {
			sendersid = bcm.fetchModifiedMap(tx.Recipient)
		}
		sendToShard[sendersid] = append(sendToShard[sendersid], tx)
	}
}

func (bcm *BrokerCommitteeMod_b2e) MsgSendingControl() {
	txfile, err := os.Open(bcm.csvPath)
	if err != nil {
		log.Panic(err)
	}
	defer txfile.Close()
	reader := csv.NewReader(txfile)
	txlist := make([]*core.Transaction, 0) // save the txs in this epoch (round)

	recoderNum := 0
	oldNum := 0

	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		if tx, ok := data2tx(data, uint64(bcm.nowDataNum)); ok {
			txlist = append(txlist, tx)
			bcm.nowDataNum++
			bcm.dataTxNums++
		} else {
			continue
		}

		// batch sending condition
		if len(txlist) == int(bcm.batchDataNum) || bcm.dataTxNums == bcm.dataTotalNum {
			// wait until all txs from the previous batch are confirmed on-chain
			bcm.b2eLock.Lock()
			itx := bcm.dealTxByBroker(txlist)
			bcm.txSending(itx)

			// Retry unallocated CTXs: as Tx1 confirms on-chain, broker balance is
			// replenished; re-run B2E on restBrokerRawMegPool until it drains.
			oldPoolSize := -1
			stuckCount := 0
			for len(bcm.restBrokerRawMegPool) != 0 {
				if len(bcm.restBrokerRawMegPool) == oldPoolSize {
					stuckCount++
					if stuckCount >= 5 {
						break // no progress for 5 s, give up
					}
				} else {
					stuckCount = 0
				}
				oldPoolSize = len(bcm.restBrokerRawMegPool)
				time.Sleep(time.Second)
				bcm.dealTxByBroker(make([]*core.Transaction, 0))
			}

			txlist = make([]*core.Transaction, 0)
			bcm.Ss.StopGap_Reset()
		}

		if bcm.dataTxNums == bcm.dataTotalNum {
			for len(bcm.restBrokerRawMegPool) != 0 {
				if len(bcm.restBrokerRawMegPool) == oldNum {
					recoderNum++
				} else {
					recoderNum = 0
				}
				bcm.dealTxByBroker(txlist)
				println("brokerTx value is ", bcm.restBrokerRawMegPool[0].Tx.Value.String())
				time.Sleep(time.Second)
				oldNum = len(bcm.restBrokerRawMegPool)
				if recoderNum >= 5 {
					break
				}
			}
			break
		}
	}

}

func (bcm *BrokerCommitteeMod_b2e) HandleBlockInfo(b *message.BlockInfoMsg) {
	bcm.sl.Slog.Printf("received from shard %d in epoch %d.\n", b.SenderShardID, b.Epoch)
	if b.BlockBodyLength == 0 {
		return
	}

	// add createConfirm
	txs := make([]*core.Transaction, 0)
	txs = append(txs, b.Broker1Txs...)
	txs = append(txs, b.Broker2Txs...)
	bcm.brokerModuleLock.Lock()
	// when accept ctx1, update all accounts
	bcm.brokerBalanceLock.Lock()
	println("block length is ", len(b.ExcutedTxs))
	for _, tx := range b.Broker1Txs {
		brokeraddress, sSid, rSid := tx.Recipient, bcm.fetchModifiedMap(tx.OriginalSender), bcm.fetchModifiedMap(tx.FinalRecipient)

		bcm.broker.LockBalance[brokeraddress][rSid].Sub(bcm.broker.LockBalance[brokeraddress][rSid], tx.Value)
		bcm.broker.BrokerBalance[brokeraddress][sSid].Add(bcm.broker.BrokerBalance[brokeraddress][sSid], tx.Value)

		fee := new(big.Float).SetInt64(tx.Fee.Int64())

		fee = fee.Mul(fee, bcm.broker.Brokerage)

		bcm.broker.ProfitBalance[brokeraddress][sSid].Add(bcm.broker.ProfitBalance[brokeraddress][sSid], fee)

	}
	bcm.add_result()

	// accumulate per-epoch stats
	bcm.epochBroker1TxNum += int(b.Broker1TxNum)
	bcm.epochBroker2TxNum += int(b.Broker2TxNum)
	bcm.epochInnerTxNum += len(b.ExcutedTxs)
	bcm.epochBlockCount++

	// accumulate confirmed tx count; unlock next batch when current batch is fully on-chain.
	// ExcutedTxs = inner txs only (RawTxHash==nil, non-allocated)
	// Broker1/2  = Tx1+Tx2 broker relay pairs; each CTX pair counts as (0.5+0.5)=1
	bcm.collectTransactions += float64(len(b.ExcutedTxs)) +
		float64(b.Broker1TxNum+b.Broker2TxNum)/2
	if bcm.collectTransactions >= float64(params.InjectSpeed) {
		bcm.collectTransactions = 0
		bcm.recordEpochStats()
		bcm.resetEpochStats()
		bcm.totalB2EIterations++
		bcm.b2eLock.Unlock()
	}

	bcm.brokerBalanceLock.Unlock()
	bcm.brokerModuleLock.Unlock()
	bcm.createConfirm(txs)
}

func (bcm *BrokerCommitteeMod_b2e) createConfirm(txs []*core.Transaction) {
	confirm1s := make([]*message.Mag1Confirm, 0)
	confirm2s := make([]*message.Mag2Confirm, 0)
	bcm.brokerModuleLock.Lock()
	for _, tx := range txs {
		if confirm1, ok := bcm.brokerConfirm1Pool[string(tx.TxHash)]; ok {
			confirm1s = append(confirm1s, confirm1)
		}
		if confirm2, ok := bcm.brokerConfirm2Pool[string(tx.TxHash)]; ok {
			confirm2s = append(confirm2s, confirm2)
		}
	}
	bcm.brokerModuleLock.Unlock()

	if len(confirm1s) != 0 {
		bcm.handleTx1ConfirmMag(confirm1s)
	}

	if len(confirm2s) != 0 {
		bcm.handleTx2ConfirmMag(confirm2s)
	}
}

func (bcm *BrokerCommitteeMod_b2e) dealTxByBroker(txs []*core.Transaction) (itxs []*core.Transaction) {
	itxs = make([]*core.Transaction, 0)
	brokerRawMegs := make([]*message.BrokerRawMeg, 0)
	//copy(brokerRawMegs, bcm.restBrokerRawMegPool)
	for _, item := range bcm.restBrokerRawMegPool {
		brokerRawMegs = append(brokerRawMegs, item)
	}
	bcm.restBrokerRawMegPool = make([]*message.BrokerRawMeg, 0)

	println("0brokerSize ", len(brokerRawMegs))
	for _, tx := range txs {
		rSid := bcm.fetchModifiedMap(tx.Recipient)
		sSid := bcm.fetchModifiedMap(tx.Sender)
		if rSid != sSid && !bcm.broker.IsBroker(tx.Recipient) && !bcm.broker.IsBroker(tx.Sender) {
			brokerRawMeg := &message.BrokerRawMeg{
				Tx:     tx,
				Broker: bcm.broker.BrokerAddress[0],
			}
			brokerRawMegs = append(brokerRawMegs, brokerRawMeg)
		} else {
			if bcm.broker.IsBroker(tx.Recipient) || bcm.broker.IsBroker(tx.Sender) {
				tx.HasBroker = true
				tx.SenderIsBroker = bcm.broker.IsBroker(tx.Sender)
			}
			itxs = append(itxs, tx)
		}
	}

	bcm.brokerBalanceLock.Lock()
	println("1brokerSize ", len(brokerRawMegs))
	alloctedBrokerRawMegs, restBrokerRawMeg := Broker2Earn.B2E(brokerRawMegs, bcm.broker.BrokerBalance)
	bcm.epochAllocatedTx += len(alloctedBrokerRawMegs)
	bcm.epochUnallocatedTx += len(restBrokerRawMeg)
	for _, item := range restBrokerRawMeg {
		bcm.restBrokerRawMegPool = append(bcm.restBrokerRawMegPool, item)
	}
	println("2brokerSize ", len(alloctedBrokerRawMegs))
	bcm.brokerBalanceLock.Unlock()
	if len(alloctedBrokerRawMegs) != 0 {
		bcm.lockToken(alloctedBrokerRawMegs)
		bcm.handleBrokerRawMag(alloctedBrokerRawMegs)
	}

	return itxs
}

func (bcm *BrokerCommitteeMod_b2e) lockToken(alloctedBrokerRawMegs []*message.BrokerRawMeg) {
	bcm.brokerBalanceLock.Lock()

	for _, brokerRawMeg := range alloctedBrokerRawMegs {
		tx := brokerRawMeg.Tx
		brokerAddress := brokerRawMeg.Broker
		rSid := bcm.fetchModifiedMap(tx.Recipient)
		bcm.broker.LockBalance[brokerAddress][rSid].Add(bcm.broker.LockBalance[brokerAddress][rSid], tx.Value)
		bcm.broker.BrokerBalance[brokerAddress][rSid].Sub(bcm.broker.BrokerBalance[brokerAddress][rSid], tx.Value)
	}

	bcm.brokerBalanceLock.Unlock()
}

func (bcm *BrokerCommitteeMod_b2e) handleBrokerType1Mes(brokerType1Megs []*message.BrokerType1Meg) {
	tx1s := make([]*core.Transaction, 0)
	for _, brokerType1Meg := range brokerType1Megs {
		ctx := brokerType1Meg.RawMeg.Tx
		tx1 := core.NewTransaction(ctx.Sender, brokerType1Meg.Broker, ctx.Value, ctx.Nonce, ctx.Fee)
		tx1.OriginalSender = ctx.Sender
		tx1.FinalRecipient = ctx.Recipient
		tx1.RawTxHash = make([]byte, len(ctx.TxHash))
		copy(tx1.RawTxHash, ctx.TxHash)
		tx1s = append(tx1s, tx1)
		confirm1 := &message.Mag1Confirm{
			RawMeg:  brokerType1Meg.RawMeg,
			Tx1Hash: tx1.TxHash,
		}
		bcm.brokerModuleLock.Lock()
		bcm.brokerConfirm1Pool[string(tx1.TxHash)] = confirm1
		bcm.brokerModuleLock.Unlock()
	}
	bcm.txSending(tx1s)
	fmt.Println("BrokerType1Mes received by shard,  add brokerTx1 len ", len(tx1s))
}

func (bcm *BrokerCommitteeMod_b2e) handleBrokerType2Mes(brokerType2Megs []*message.BrokerType2Meg) {
	tx2s := make([]*core.Transaction, 0)
	for _, mes := range brokerType2Megs {
		ctx := mes.RawMeg.Tx
		tx2 := core.NewTransaction(mes.Broker, ctx.Recipient, ctx.Value, ctx.Nonce, ctx.Fee)
		tx2.OriginalSender = ctx.Sender
		tx2.FinalRecipient = ctx.Recipient
		tx2.RawTxHash = make([]byte, len(ctx.TxHash))
		copy(tx2.RawTxHash, ctx.TxHash)
		tx2s = append(tx2s, tx2)

		confirm2 := &message.Mag2Confirm{
			RawMeg:  mes.RawMeg,
			Tx2Hash: tx2.TxHash,
		}
		bcm.brokerModuleLock.Lock()
		bcm.brokerConfirm2Pool[string(tx2.TxHash)] = confirm2
		bcm.brokerModuleLock.Unlock()
	}
	bcm.txSending(tx2s)
	fmt.Println("broker tx2 add to pool len ", len(tx2s))
}

// get the digest of rawMeg
func (bcm *BrokerCommitteeMod_b2e) getBrokerRawMagDigest(r *message.BrokerRawMeg) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

func (bcm *BrokerCommitteeMod_b2e) handleBrokerRawMag(brokerRawMags []*message.BrokerRawMeg) {
	b := bcm.broker
	brokerType1Mags := make([]*message.BrokerType1Meg, 0)
	fmt.Println("broker receive ctx ", len(brokerRawMags))
	bcm.brokerModuleLock.Lock()
	for _, meg := range brokerRawMags {
		b.BrokerRawMegs[string(bcm.getBrokerRawMagDigest(meg))] = meg
		brokerType1Mag := &message.BrokerType1Meg{
			RawMeg:   meg,
			Hcurrent: 0,
			Broker:   meg.Broker,
		}
		brokerType1Mags = append(brokerType1Mags, brokerType1Mag)
	}
	bcm.brokerModuleLock.Unlock()
	bcm.handleBrokerType1Mes(brokerType1Mags)
}

func (bcm *BrokerCommitteeMod_b2e) handleTx1ConfirmMag(mag1confirms []*message.Mag1Confirm) {
	brokerType2Mags := make([]*message.BrokerType2Meg, 0)
	b := bcm.broker

	fmt.Println("receive confirm  brokerTx1 len ", len(mag1confirms))
	bcm.brokerModuleLock.Lock()
	for _, mag1confirm := range mag1confirms {
		RawMeg := mag1confirm.RawMeg
		_, ok := b.BrokerRawMegs[string(bcm.getBrokerRawMagDigest(RawMeg))]
		if !ok {
			fmt.Println("raw message is not exited,tx1 confirms failure !")
			continue
		}
		b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)] = append(b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)], string(mag1confirm.Tx1Hash))
		brokerType2Mag := &message.BrokerType2Meg{
			Broker: RawMeg.Broker,
			RawMeg: RawMeg,
		}
		brokerType2Mags = append(brokerType2Mags, brokerType2Mag)
	}
	bcm.brokerModuleLock.Unlock()
	bcm.handleBrokerType2Mes(brokerType2Mags)
}

func (bcm *BrokerCommitteeMod_b2e) handleTx2ConfirmMag(mag2confirms []*message.Mag2Confirm) {
	b := bcm.broker
	fmt.Println("receive confirm  brokerTx2 len ", len(mag2confirms))
	num := 0
	bcm.brokerModuleLock.Lock()
	for _, mag2confirm := range mag2confirms {
		RawMeg := mag2confirm.RawMeg
		b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)] = append(b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)], string(mag2confirm.Tx2Hash))
		if len(b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)]) == 2 {
			num++
		} else {
			fmt.Println(len(b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)]))
		}
	}
	bcm.brokerModuleLock.Unlock()
	fmt.Println("finish ctx with adding tx1 and tx2 to txpool,len", num)
}

func (bcm *BrokerCommitteeMod_b2e) add_result() {

	for brokerAddress, shardMap := range bcm.broker.BrokerBalance {
		a := ""
		b := ""
		c := ""
		for shardId, balance := range shardMap {
			a += balance.String() + ","
			b += bcm.broker.LockBalance[brokerAddress][shardId].String() + ","
			c += bcm.broker.ProfitBalance[brokerAddress][shardId].String() + ","
		}
		a += "\n"
		b += "\n"
		c += "\n"
		bcm.Result_lockBalance[brokerAddress] = append(bcm.Result_lockBalance[brokerAddress], b)
		bcm.Result_brokerBalance[brokerAddress] = append(bcm.Result_brokerBalance[brokerAddress], a)
		bcm.Result_Profit[brokerAddress] = append(bcm.Result_Profit[brokerAddress], c)
	}
}
func (bcm *BrokerCommitteeMod_b2e) Result_save() {

	// write to .csv file
	dirpath := params.DataWrite_path + "brokerRsult/"
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	for brokerAddress, _ := range bcm.broker.BrokerBalance {
		targetPath0 := dirpath + brokerAddress + "_lockBalance.csv"
		targetPath1 := dirpath + brokerAddress + "_brokerBalance.csv"
		targetPath2 := dirpath + brokerAddress + "_Profit.csv"
		bcm.Wirte_result(targetPath0, bcm.Result_lockBalance[brokerAddress])
		bcm.Wirte_result(targetPath1, bcm.Result_brokerBalance[brokerAddress])
		bcm.Wirte_result(targetPath2, bcm.Result_Profit[brokerAddress])
	}
}
func (bcm *BrokerCommitteeMod_b2e) recordEpochStats() {
	dirpath := params.DataWrite_path + "epoch_stats/"
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}

	// compute cumulative total profit in ETH (sum across all brokers and shards)
	// ProfitBalance stores fee * Brokerage as *big.Float, in Wei units
	totalProfit := new(big.Float).SetFloat64(0)
	for _, shardMap := range bcm.broker.ProfitBalance {
		for _, profit := range shardMap {
			totalProfit.Add(totalProfit, profit)
		}
	}
	// convert Wei to ETH (1 ETH = 1e18 Wei)
	totalProfit.Quo(totalProfit, new(big.Float).SetFloat64(1e18))
	profitETH, _ := totalProfit.Float64()

	targetPath := dirpath + fmt.Sprintf("epoch_stats_Broker2Earn_inject%d_shard%d.csv", params.InjectSpeed, params.ShardNum)

	// write header if file does not yet exist
	_, statErr := os.Stat(targetPath)
	file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	if os.IsNotExist(statErr) {
		w.Write([]string{
			"Epoch", "Injected_Tx",
			"B2E_Allocated_Tx", "B2E_Unallocated_Tx",
			"Broker_Service_Tx", "Inner_Tx",
			"Total_Profit_ETH", "Active_Broker_Count", "Block_Count",
		})
	}
	w.Write([]string{
		strconv.Itoa(bcm.totalB2EIterations),
		strconv.Itoa(params.InjectSpeed),
		strconv.Itoa(bcm.epochAllocatedTx),
		strconv.Itoa(bcm.epochUnallocatedTx),
		strconv.Itoa(bcm.epochBroker1TxNum + bcm.epochBroker2TxNum),
		strconv.Itoa(bcm.epochInnerTxNum),
		fmt.Sprintf("%.6f", profitETH),
		strconv.Itoa(len(bcm.broker.BrokerAddress)),
		strconv.Itoa(bcm.epochBlockCount),
	})
	w.Flush()
}

func (bcm *BrokerCommitteeMod_b2e) resetEpochStats() {
	bcm.epochBroker1TxNum = 0
	bcm.epochBroker2TxNum = 0
	bcm.epochInnerTxNum = 0
	bcm.epochBlockCount = 0
	bcm.epochAllocatedTx = 0
	bcm.epochUnallocatedTx = 0
}

func (bcm *BrokerCommitteeMod_b2e) Wirte_result(targetPath string, resultStr []string) {

	f, err := os.Open(targetPath)
	if err != nil && os.IsNotExist(err) {
		file, er := os.Create(targetPath)
		if er != nil {
			panic(er)
		}
		defer file.Close()

		w := csv.NewWriter(file)
		w.Flush()
		for _, str := range resultStr {
			str_arry := strings.Split(str, ",")
			w.Write(str_arry[0 : len(str_arry)-1])
			w.Flush()
		}
	} else {
		file, err := os.OpenFile(targetPath, os.O_APPEND|os.O_RDWR, 0666)

		if err != nil {
			log.Panic(err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		err = writer.Write(resultStr)
		if err != nil {
			log.Panic()
		}
		writer.Flush()
	}
	f.Close()
}
