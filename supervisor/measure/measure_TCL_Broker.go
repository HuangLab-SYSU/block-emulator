package measure

import (
	"blockEmulator/message"
	"fmt"
	"strconv"
	"time"
)

// to test average Transaction_Confirm_Latency (TCL) in this system
type TestModule_TCL_Broker struct {
	epochID int

	totTxLatencyEpoch     []float64 // record the Transaction_Confirm_Latency in each epoch
	broker1CommitLatency  []int64
	broker2CommitLatency  []int64
	ctxCommitLatency      []int64
	normalTxCommitLatency []int64

	normalTxNum  []int
	broker1TxNum []int
	broker2TxNum []int
	txNum        []float64 // record the txNumber in each epoch

	brokerTxMap map[string]time.Time // map: origin raw tx' hash to the time when the corresponding broker1 tx was added into the pool.
}

func NewTestModule_TCL_Broker() *TestModule_TCL_Broker {
	return &TestModule_TCL_Broker{
		epochID:               -1,
		totTxLatencyEpoch:     make([]float64, 0),
		broker1CommitLatency:  make([]int64, 0),
		broker2CommitLatency:  make([]int64, 0),
		ctxCommitLatency:      make([]int64, 0),
		normalTxCommitLatency: make([]int64, 0),

		txNum:        make([]float64, 0),
		normalTxNum:  make([]int, 0),
		broker1TxNum: make([]int, 0),
		broker2TxNum: make([]int, 0),
		brokerTxMap:  make(map[string]time.Time),
	}
}

func (tml *TestModule_TCL_Broker) OutputMetricName() string {
	return "Transaction_Confirm_Latency"
}

// modified TCL
func (tml *TestModule_TCL_Broker) UpdateMeasureRecord(b *message.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}

	epochid := b.Epoch

	// extend
	for tml.epochID < epochid {
		tml.txNum = append(tml.txNum, 0)
		tml.totTxLatencyEpoch = append(tml.totTxLatencyEpoch, 0)

		tml.broker1CommitLatency = append(tml.broker1CommitLatency, 0)
		tml.broker2CommitLatency = append(tml.broker2CommitLatency, 0)
		tml.normalTxCommitLatency = append(tml.normalTxCommitLatency, 0)
		tml.ctxCommitLatency = append(tml.ctxCommitLatency, 0)

		tml.broker1TxNum = append(tml.broker1TxNum, 0)
		tml.broker2TxNum = append(tml.broker2TxNum, 0)
		tml.normalTxNum = append(tml.normalTxNum, 0)

		tml.epochID++
	}

	tml.broker1TxNum[epochid] += len(b.Broker1Txs)
	tml.broker2TxNum[epochid] += len(b.Broker2Txs)
	tml.normalTxNum[epochid] += len(b.InnerShardTxs)
	tml.txNum[epochid] += float64(len(b.InnerShardTxs)) + float64(len(b.Broker1Txs)+len(b.Broker2Txs))/2

	// normal txs
	for _, tx := range b.InnerShardTxs {
		tml.totTxLatencyEpoch[epochid] += b.CommitTime.Sub(tx.Time).Seconds()
		tml.normalTxCommitLatency[epochid] += int64(b.CommitTime.Sub(tx.Time).Milliseconds())
	}
	// broker
	for _, b1tx := range b.Broker1Txs {
		tml.brokerTxMap[string(b1tx.RawTxHash)] = b1tx.Time
		tml.broker1CommitLatency[epochid] += int64(b.CommitTime.Sub(b1tx.Time).Milliseconds())
	}
	for _, b2tx := range b.Broker2Txs {
		if b1txProposeTime, ok := tml.brokerTxMap[string(b2tx.RawTxHash)]; ok {
			tml.totTxLatencyEpoch[epochid] += b.CommitTime.Sub(b1txProposeTime).Seconds()
			tml.ctxCommitLatency[epochid] += b.CommitTime.Sub(b1txProposeTime).Milliseconds()
		} else {
			fmt.Println("Missing a broker1 tx. ")
		}
		tml.broker2CommitLatency[epochid] += int64(b.CommitTime.Sub(b2tx.Time).Milliseconds())
	}
}

func (tml *TestModule_TCL_Broker) HandleExtraMessage(msg []byte) {}

func (tml *TestModule_TCL_Broker) OutputRecord() (perEpochLatency []float64, totLatency float64) {
	tml.writeToCSV()

	// calculate the simple result
	perEpochLatency = make([]float64, 0)
	latencySum := 0.0
	totTxNum := 0.0

	for eid, totLatency := range tml.totTxLatencyEpoch {
		perEpochLatency = append(perEpochLatency, totLatency/tml.txNum[eid])
		latencySum += totLatency
		totTxNum += tml.txNum[eid]
	}
	totLatency = latencySum / totTxNum
	return
}

func (tml *TestModule_TCL_Broker) writeToCSV() {
	fileName := tml.OutputMetricName()
	measureName := []string{"EpochID",
		"Total tx # in this epoch",
		"Normal tx # in this epoch",
		"Broker1 tx # in this epoch",
		"Broker2 tx # in this epoch",
		"Sum of Broker1 TCL (ms) (Duration: Broker1 Tx Propose -> Broker1 Tx Commit)",
		"Sum of Broker2 TCL (ms) (Duration: Broker2 Tx Propose -> Broker2 Tx Commit)",
		"Sum of innerShardTx TCL (ms)",
		"Sum of CTX TCL (ms) (Duration: Broker1 Tx Propose -> Broker2 Tx Commit)",
		"Sum of All Tx TCL (sec.)"}

	measureVals := make([][]string, 0)
	for eid, totTxInE := range tml.txNum {
		csvLine := []string{
			strconv.Itoa(eid),
			strconv.FormatFloat(totTxInE, 'f', '8', 64),
			strconv.Itoa(tml.normalTxNum[eid]),
			strconv.Itoa(tml.broker1TxNum[eid]),
			strconv.Itoa(tml.broker2TxNum[eid]),
			strconv.FormatInt(tml.broker1CommitLatency[eid], 10),
			strconv.FormatInt(tml.broker2CommitLatency[eid], 10),
			strconv.FormatInt(tml.normalTxCommitLatency[eid], 10),
			strconv.FormatInt(tml.ctxCommitLatency[eid], 10),
			strconv.FormatFloat(tml.totTxLatencyEpoch[eid], 'f', '8', 64),
		}
		measureVals = append(measureVals, csvLine)
	}
	WriteMetricsToCSV(fileName, measureName, measureVals)
}
