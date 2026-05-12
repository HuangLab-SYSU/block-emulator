package measure

import (
	"blockEmulator/message"
	"time"
)

// to test average Transaction_Confirm_Latency (TCL) in this system
type TestModule_TCL_Broker struct {
	epochID           int
	totTxLatencyEpoch []float64            // record the Transaction_Confirm_Latency in each epoch
	txNum             []float64            // record the txNumber in each epoch
	brokerTxMap       map[string]time.Time // map: origin raw tx' hash to the time when the corresponding broker1 tx was added into the pool.
}

func NewTestModule_TCL_Broker() *TestModule_TCL_Broker {
	return &TestModule_TCL_Broker{
		epochID:           -1,
		totTxLatencyEpoch: make([]float64, 0),
		txNum:             make([]float64, 0),
		brokerTxMap:       make(map[string]time.Time),
	}
}

func (tml *TestModule_TCL_Broker) OutputMetricName() string {
	return "Transaction_Confirm_Latency"
}

func (tml *TestModule_TCL_Broker) OutputMetricTitle() string {
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
		tml.epochID++
	}

	// calculate txs
	for _, tx := range b.ExcutedTxs {
		if !tx.Time.IsZero() {
			tml.totTxLatencyEpoch[epochid] += b.CommitTime.Sub(tx.Time).Seconds()
			tml.txNum[epochid]++
		}
	}
	// broker
	for _, b1tx := range b.Broker1Txs {
		tml.brokerTxMap[string(b1tx.RawTxHash)] = b1tx.Time
	}
	for _, b2tx := range b.Broker2Txs {
		if t, ok := tml.brokerTxMap[string(b2tx.RawTxHash)]; ok {
			tml.totTxLatencyEpoch[epochid] += b.ProposeTime.Sub(t).Seconds()
		}
	}
}

func (tml *TestModule_TCL_Broker) HandleExtraMessage(msg []byte) {}

func (tml *TestModule_TCL_Broker) OutputRecord() (perEpochLatency []float64, totLatency float64) {
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
