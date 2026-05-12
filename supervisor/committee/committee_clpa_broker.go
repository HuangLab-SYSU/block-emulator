package committee

import (
	"blockEmulator/broker"
	"blockEmulator/core"
	"blockEmulator/message"
	"blockEmulator/networks"
	"blockEmulator/params"
	"blockEmulator/partition"
	"blockEmulator/supervisor/signal"
	"blockEmulator/supervisor/supervisor_log"
	"blockEmulator/utils"
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// CLPA committee operations
type CLPACommitteeMod_Broker struct {
	csvPath      string
	dataTotalNum int
	nowDataNum   int
	batchDataNum int

	// additional variants
	clpaLock            sync.Mutex
	clpaGraph           *partition.CLPAState
	modifiedMap         map[string]uint64
	clpaLastRunningTime time.Time
	clpaFreq            int

	//broker related  attributes avatar
	broker             *broker.Broker
	brokerConfirm1Pool map[string]*message.Mag1Confirm
	brokerConfirm2Pool map[string]*message.Mag2Confirm
	brokerTxPool       []*core.Transaction
	brokerModuleLock   sync.Mutex

	// logger module
	sl *supervisor_log.SupervisorLog

	// control components
	Ss          *signal.StopSignal // to control the stop message sending
	IpNodeTable map[uint64]map[uint64]string
}

func NewCLPACommitteeMod_Broker(Ip_nodeTable map[uint64]map[uint64]string, Ss *signal.StopSignal, sl *supervisor_log.SupervisorLog, csvFilePath string, dataNum, batchNum, clpaFrequency int) *CLPACommitteeMod_Broker {
	cg := new(partition.CLPAState)
	cg.Init_CLPAState(0.5, 100, params.ShardNum)

	broker := new(broker.Broker)
	broker.NewBroker(nil)

	return &CLPACommitteeMod_Broker{
		csvPath:             csvFilePath,
		dataTotalNum:        dataNum,
		batchDataNum:        batchNum,
		nowDataNum:          0,
		clpaGraph:           cg,
		modifiedMap:         make(map[string]uint64),
		clpaFreq:            clpaFrequency,
		clpaLastRunningTime: time.Time{},
		brokerConfirm1Pool:  make(map[string]*message.Mag1Confirm),
		brokerConfirm2Pool:  make(map[string]*message.Mag2Confirm),
		brokerTxPool:        make([]*core.Transaction, 0),
		broker:              broker,
		IpNodeTable:         Ip_nodeTable,
		Ss:                  Ss,
		sl:                  sl,
	}
}

// for CLPA_Broker committee, it only handle the extra CInner2CrossTx message.
func (ccm *CLPACommitteeMod_Broker) HandleOtherMessage(msg []byte) {
	msgType, content := message.SplitMessage(msg)
	if msgType != message.CInner2CrossTx {
		return
	}
	itct := new(message.InnerTx2CrossTx)
	err := json.Unmarshal(content, itct)
	if err != nil {
		log.Panic()
	}
	itxs := ccm.dealTxByBroker(itct.Txs)
	ccm.txSending(itxs)
}

func (ccm *CLPACommitteeMod_Broker) fetchModifiedMap(key string) uint64 {
	if val, ok := ccm.modifiedMap[key]; !ok {
		return uint64(utils.Addr2Shard(key))
	} else {
		return val
	}
}

func (ccm *CLPACommitteeMod_Broker) txSending(txlist []*core.Transaction) {
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
				go networks.TcpDial(send_msg, ccm.IpNodeTable[sid][0])
			}
			sendToShard = make(map[uint64][]*core.Transaction)
			time.Sleep(time.Second)
		}
		if idx == len(txlist) {
			break
		}
		tx := txlist[idx]
		ccm.clpaLock.Lock()
		sendersid := ccm.fetchModifiedMap(tx.Sender)

		if ccm.broker.IsBroker(tx.Sender) {
			sendersid = ccm.fetchModifiedMap(tx.Recipient)
		}

		ccm.clpaLock.Unlock()
		sendToShard[sendersid] = append(sendToShard[sendersid], tx)
	}
}

func (ccm *CLPACommitteeMod_Broker) MsgSendingControl() {
	txfile, err := os.Open(ccm.csvPath)
	if err != nil {
		log.Panic(err)
	}
	defer txfile.Close()
	reader := csv.NewReader(txfile)
	txlist := make([]*core.Transaction, 0) // save the txs in this epoch (round)

	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		if tx, ok := data2tx(data, uint64(ccm.nowDataNum)); ok {
			txlist = append(txlist, tx)
			ccm.nowDataNum++
		} else {
			continue
		}

		// batch sending condition
		if len(txlist) == int(ccm.batchDataNum) || ccm.nowDataNum == ccm.dataTotalNum {
			// set the algorithm timer begins
			if ccm.clpaLastRunningTime.IsZero() {
				ccm.clpaLastRunningTime = time.Now()
			}

			itx := ccm.dealTxByBroker(txlist)

			ccm.txSending(itx)

			// reset the variants about tx sending
			txlist = make([]*core.Transaction, 0)
			ccm.Ss.StopGap_Reset()
		}

		if !ccm.clpaLastRunningTime.IsZero() && time.Since(ccm.clpaLastRunningTime) >= time.Duration(ccm.clpaFreq)*time.Second {
			ccm.clpaLock.Lock()
			mmap, _ := ccm.clpaGraph.CLPA_Partition()
			ccm.clpaMapSend(mmap)
			for key, val := range mmap {
				ccm.modifiedMap[key] = val
			}
			ccm.clpaReset()
			ccm.clpaLock.Unlock()
			time.Sleep(10 * time.Second)
			ccm.clpaLastRunningTime = time.Now()
		}

		if ccm.nowDataNum == ccm.dataTotalNum {
			break
		}
	}

	// all transactions are sent. keep sending partition message...
	for !ccm.Ss.GapEnough() { // wait all txs to be handled
		time.Sleep(time.Second)
		if time.Since(ccm.clpaLastRunningTime) >= time.Duration(ccm.clpaFreq)*time.Second {
			ccm.clpaLock.Lock()
			mmap, _ := ccm.clpaGraph.CLPA_Partition()
			ccm.clpaMapSend(mmap)
			for key, val := range mmap {
				ccm.modifiedMap[key] = val
			}
			ccm.clpaReset()
			ccm.clpaLock.Unlock()
			time.Sleep(10 * time.Second)
			ccm.clpaLastRunningTime = time.Now()
		}
	}
}

func (ccm *CLPACommitteeMod_Broker) clpaMapSend(m map[string]uint64) {
	// send partition modified Map message
	pm := message.PartitionModifiedMap{
		PartitionModified: m,
	}
	pmByte, err := json.Marshal(pm)
	if err != nil {
		log.Panic()
	}
	send_msg := message.MergeMessage(message.CPartitionMsg, pmByte)
	// send to worker shards
	for i := uint64(0); i < uint64(params.ShardNum); i++ {
		networks.TcpDial(send_msg, ccm.IpNodeTable[i][0])
	}
	ccm.sl.Slog.Println("Supervisor: all partition map message has been sent. ")
}

func (ccm *CLPACommitteeMod_Broker) clpaReset() {
	ccm.clpaGraph = new(partition.CLPAState)
	ccm.clpaGraph.Init_CLPAState(0.5, 100, params.ShardNum)
	for key, val := range ccm.modifiedMap {
		ccm.clpaGraph.PartitionMap[partition.Vertex{Addr: key}] = int(val)
	}
}

func (ccm *CLPACommitteeMod_Broker) HandleBlockInfo(b *message.BlockInfoMsg) {
	ccm.sl.Slog.Printf("received from shard %d in epoch %d.\n", b.SenderShardID, b.Epoch)
	if b.BlockBodyLength == 0 {
		return
	}

	// add createConfirm
	txs := make([]*core.Transaction, 0)
	txs = append(txs, b.Broker1Txs...)
	txs = append(txs, b.Broker2Txs...)
	ccm.createConfirm(txs)

	ccm.clpaLock.Lock()
	for _, tx := range b.ExcutedTxs {
		if tx.HasBroker {
			continue
		}
		ccm.clpaGraph.AddEdge(partition.Vertex{Addr: tx.Sender}, partition.Vertex{Addr: tx.Recipient})
	}
	for _, b1tx := range b.Broker1Txs {
		ccm.clpaGraph.AddEdge(partition.Vertex{Addr: b1tx.OriginalSender}, partition.Vertex{Addr: b1tx.FinalRecipient})
	}
	ccm.clpaLock.Unlock()
}

func (ccm *CLPACommitteeMod_Broker) createConfirm(txs []*core.Transaction) {
	confirm1s := make([]*message.Mag1Confirm, 0)
	confirm2s := make([]*message.Mag2Confirm, 0)
	ccm.brokerModuleLock.Lock()
	for _, tx := range txs {
		if confirm1, ok := ccm.brokerConfirm1Pool[string(tx.TxHash)]; ok {
			confirm1s = append(confirm1s, confirm1)
		}
		if confirm2, ok := ccm.brokerConfirm2Pool[string(tx.TxHash)]; ok {
			confirm2s = append(confirm2s, confirm2)
		}
	}
	ccm.brokerModuleLock.Unlock()

	if len(confirm1s) != 0 {
		ccm.handleTx1ConfirmMag(confirm1s)
	}

	if len(confirm2s) != 0 {
		ccm.handleTx2ConfirmMag(confirm2s)
	}
}

func (ccm *CLPACommitteeMod_Broker) dealTxByBroker(txs []*core.Transaction) (itxs []*core.Transaction) {
	itxs = make([]*core.Transaction, 0)
	brokerRawMegs := make([]*message.BrokerRawMeg, 0)
	for _, tx := range txs {
		ccm.clpaLock.Lock()
		rSid := ccm.fetchModifiedMap(tx.Recipient)
		sSid := ccm.fetchModifiedMap(tx.Sender)
		ccm.clpaLock.Unlock()
		if rSid != sSid && !ccm.broker.IsBroker(tx.Recipient) && !ccm.broker.IsBroker(tx.Sender) {
			brokerRawMeg := &message.BrokerRawMeg{
				Tx:     tx,
				Broker: ccm.broker.BrokerAddress[0],
			}
			brokerRawMegs = append(brokerRawMegs, brokerRawMeg)
		} else {
			if ccm.broker.IsBroker(tx.Recipient) || ccm.broker.IsBroker(tx.Sender) {
				tx.HasBroker = true
				tx.SenderIsBroker = ccm.broker.IsBroker(tx.Sender)
			}
			itxs = append(itxs, tx)
		}
	}
	if len(brokerRawMegs) != 0 {
		ccm.handleBrokerRawMag(brokerRawMegs)
	}
	return itxs
}

func (ccm *CLPACommitteeMod_Broker) handleBrokerType1Mes(brokerType1Megs []*message.BrokerType1Meg) {
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
		ccm.brokerModuleLock.Lock()
		ccm.brokerConfirm1Pool[string(tx1.TxHash)] = confirm1
		ccm.brokerModuleLock.Unlock()
	}
	ccm.txSending(tx1s)
	fmt.Println("BrokerType1Mes received by shard,  add brokerTx1 len ", len(tx1s))
}

func (ccm *CLPACommitteeMod_Broker) handleBrokerType2Mes(brokerType2Megs []*message.BrokerType2Meg) {
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
		ccm.brokerModuleLock.Lock()
		ccm.brokerConfirm2Pool[string(tx2.TxHash)] = confirm2
		ccm.brokerModuleLock.Unlock()
	}
	ccm.txSending(tx2s)
	fmt.Println("broker tx2 add to pool len ", len(tx2s))
}

// get the digest of rawMeg
func (ccm *CLPACommitteeMod_Broker) getBrokerRawMagDigest(r *message.BrokerRawMeg) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

func (ccm *CLPACommitteeMod_Broker) handleBrokerRawMag(brokerRawMags []*message.BrokerRawMeg) {
	b := ccm.broker
	brokerType1Mags := make([]*message.BrokerType1Meg, 0)
	fmt.Println("broker receive ctx ", len(brokerRawMags))
	ccm.brokerModuleLock.Lock()
	for _, meg := range brokerRawMags {
		b.BrokerRawMegs[string(ccm.getBrokerRawMagDigest(meg))] = meg

		brokerType1Mag := &message.BrokerType1Meg{
			RawMeg:   meg,
			Hcurrent: 0,
			Broker:   meg.Broker,
		}
		brokerType1Mags = append(brokerType1Mags, brokerType1Mag)
	}
	ccm.brokerModuleLock.Unlock()
	ccm.handleBrokerType1Mes(brokerType1Mags)
}

func (ccm *CLPACommitteeMod_Broker) handleTx1ConfirmMag(mag1confirms []*message.Mag1Confirm) {
	brokerType2Mags := make([]*message.BrokerType2Meg, 0)
	b := ccm.broker

	fmt.Println("receive confirm  brokerTx1 len ", len(mag1confirms))
	ccm.brokerModuleLock.Lock()
	for _, mag1confirm := range mag1confirms {
		RawMeg := mag1confirm.RawMeg
		_, ok := b.BrokerRawMegs[string(ccm.getBrokerRawMagDigest(RawMeg))]
		if !ok {
			fmt.Println("raw message is not exited,tx1 confirms failure !")
			continue
		}
		b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)] = append(b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)], string(mag1confirm.Tx1Hash))

		brokerType2Mag := &message.BrokerType2Meg{
			Broker: ccm.broker.BrokerAddress[0],
			RawMeg: RawMeg,
		}
		brokerType2Mags = append(brokerType2Mags, brokerType2Mag)
	}
	ccm.brokerModuleLock.Unlock()
	ccm.handleBrokerType2Mes(brokerType2Mags)
}

func (ccm *CLPACommitteeMod_Broker) handleTx2ConfirmMag(mag2confirms []*message.Mag2Confirm) {
	b := ccm.broker
	fmt.Println("receive confirm  brokerTx2 len ", len(mag2confirms))
	num := 0
	ccm.brokerModuleLock.Lock()
	for _, mag2confirm := range mag2confirms {
		RawMeg := mag2confirm.RawMeg
		b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)] = append(b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)], string(mag2confirm.Tx2Hash))
		if len(b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)]) == 2 {
			num++
		} else {
			fmt.Println(len(b.RawTx2BrokerTx[string(RawMeg.Tx.TxHash)]))
		}
	}
	ccm.brokerModuleLock.Unlock()
	fmt.Println("finish ctx with adding tx1 and tx2 to txpool,len", num)
}

func (bcm *CLPACommitteeMod_Broker) Result_save() {
}
