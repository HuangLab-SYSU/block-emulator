package measure

import "blockEmulator/message"

// to test cross-transaction rate
type TestCrossTxRate_Relay struct {
	epochID       int
	totTxNum      []float64
	totCrossTxNum []float64
}

func NewTestCrossTxRate_Relay() *TestCrossTxRate_Relay {
	return &TestCrossTxRate_Relay{
		epochID:       -1,
		totTxNum:      make([]float64, 0),
		totCrossTxNum: make([]float64, 0),
	}
}

func (tctr *TestCrossTxRate_Relay) OutputMetricName() string {
	return "CrossTransaction_ratio"
}

func (tctr *TestCrossTxRate_Relay) UpdateMeasureRecord(b *message.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}

	epochid := b.Epoch
	r1TxNum := len(b.Relay1Txs)
	r2TxNum := len(b.Relay2Txs)

	// extend
	for tctr.epochID < epochid {
		tctr.totTxNum = append(tctr.totTxNum, 0)
		tctr.totCrossTxNum = append(tctr.totCrossTxNum, 0)
		tctr.epochID++
	}

	tctr.totCrossTxNum[epochid] += float64(r1TxNum+r2TxNum) / 2
	tctr.totTxNum[epochid] += float64(r1TxNum+r2TxNum)/2 + float64(len(b.InterShardTxs))
}

func (tctr *TestCrossTxRate_Relay) HandleExtraMessage([]byte) {}

func (tctr *TestCrossTxRate_Relay) OutputRecord() (perEpochCTXratio []float64, totCTXratio float64) {
	perEpochCTXratio = make([]float64, 0)
	allEpoch_totTxNum := 0.0
	allEpoch_ctxNum := 0.0
	for eid, totTxN := range tctr.totTxNum {
		perEpochCTXratio = append(perEpochCTXratio, tctr.totCrossTxNum[eid]/totTxN)
		allEpoch_totTxNum += totTxN
		allEpoch_ctxNum += tctr.totCrossTxNum[eid]
	}
	perEpochCTXratio = append(perEpochCTXratio, allEpoch_totTxNum)
	perEpochCTXratio = append(perEpochCTXratio, allEpoch_ctxNum)

	return perEpochCTXratio, allEpoch_ctxNum / allEpoch_totTxNum
}
