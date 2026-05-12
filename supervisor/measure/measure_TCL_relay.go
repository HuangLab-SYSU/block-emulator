package measure

import (
	"blockEmulator/message"
)

// to test average Transaction_Confirm_Latency (TCL)  in this system
type TestModule_TCL_Relay struct {
	epochID           int
	totTxLatencyEpoch []float64 // record the Transaction_Confirm_Latency in each epoch
	txNum             []float64 // record the txNumber in each epoch
}

func NewTestModule_TCL_Relay() *TestModule_TCL_Relay {
	return &TestModule_TCL_Relay{
		epochID:           -1,
		totTxLatencyEpoch: make([]float64, 0),
		txNum:             make([]float64, 0),
	}
}

func (tml *TestModule_TCL_Relay) OutputMetricName() string {
	return "Transaction_Confirm_Latency"
}

func (tml *TestModule_TCL_Relay) OutputMetricTitle() string {
	return "Transaction_Confirm_Latency"
}

// modified latency
func (tml *TestModule_TCL_Relay) UpdateMeasureRecord(b *message.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}

	epochid := b.Epoch
	txs := b.ExcutedTxs
	mTime := b.CommitTime

	// extend
	for tml.epochID < epochid {
		tml.txNum = append(tml.txNum, 0)
		tml.totTxLatencyEpoch = append(tml.totTxLatencyEpoch, 0)
		tml.epochID++
	}

	for _, tx := range txs {
		if !tx.Time.IsZero() {
			tml.totTxLatencyEpoch[epochid] += mTime.Sub(tx.Time).Seconds()
			tml.txNum[epochid]++
		}
	}
}

func (tml *TestModule_TCL_Relay) HandleExtraMessage([]byte) {}

func (tml *TestModule_TCL_Relay) OutputRecord() (perEpochLatency []float64, totLatency float64) {
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
