package measure

import (
	"blockEmulator/message"
	"fmt"
)

// to test cross-transaction rate
type TestCrossTxRate_Broker struct {
	epochID       int
	totTxNum      []float64
	totCrossTxNum []float64
	b2num, b1num  int
}

func NewTestCrossTxRate_Broker() *TestCrossTxRate_Broker {
	return &TestCrossTxRate_Broker{
		epochID:       -1,
		totTxNum:      make([]float64, 0),
		totCrossTxNum: make([]float64, 0),
		b2num:         0,
		b1num:         0,
	}
}

func (tctr *TestCrossTxRate_Broker) OutputMetricName() string {
	return "CrossTransaction_ratio"
}
func (tctr *TestCrossTxRate_Broker) OutputMetricTitle() string {
	return "CrossTransaction_ratio"
}

func (tctr *TestCrossTxRate_Broker) UpdateMeasureRecord(b *message.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}
	epochid := b.Epoch
	// extend
	for tctr.epochID < epochid {
		tctr.totTxNum = append(tctr.totTxNum, 0)
		tctr.totCrossTxNum = append(tctr.totCrossTxNum, 0)
		tctr.epochID++
	}

	tctr.totCrossTxNum[epochid] += (float64(b.Broker1TxNum) + float64(b.Broker2TxNum)) / 2
	tctr.totTxNum[epochid] += float64(len(b.ExcutedTxs)) + (float64(b.Broker1TxNum)+float64(b.Broker2TxNum))/2
	tctr.b2num += int(b.Broker2TxNum)
	tctr.b1num += int(b.Broker1TxNum)
}

func (tctr *TestCrossTxRate_Broker) HandleExtraMessage([]byte) {}

func (tctr *TestCrossTxRate_Broker) OutputRecord() (perEpochCTXratio []float64, totCTXratio float64) {
	fmt.Println(tctr.b2num)
	fmt.Println(tctr.b1num)

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
