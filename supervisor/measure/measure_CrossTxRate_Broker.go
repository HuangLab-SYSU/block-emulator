package measure

import (
	"blockEmulator/message"
	"strconv"
)

// to test cross-transaction rate
type TestCrossTxRate_Broker struct {
	epochID       int
	totTxNum      []float64
	totCrossTxNum []float64

	broker1TxNum []int // record how many broker1 txs in an epoch.
	broker2TxNum []int // record how many broker2 txs in an epoch.
	normalTxNum  []int // record how many normal txs in an epoch.
}

func NewTestCrossTxRate_Broker() *TestCrossTxRate_Broker {
	return &TestCrossTxRate_Broker{
		epochID:       -1,
		totTxNum:      make([]float64, 0),
		totCrossTxNum: make([]float64, 0),

		broker1TxNum: make([]int, 0),
		broker2TxNum: make([]int, 0),
		normalTxNum:  make([]int, 0),
	}
}

func (tctr *TestCrossTxRate_Broker) OutputMetricName() string {
	return "CrossTransaction_ratio"
}

func (tctr *TestCrossTxRate_Broker) UpdateMeasureRecord(b *message.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}
	b1TxNum := len(b.Broker1Txs)
	b2TxNum := len(b.Broker2Txs)
	epochid := b.Epoch
	// extend
	for tctr.epochID < epochid {
		tctr.totTxNum = append(tctr.totTxNum, 0)
		tctr.totCrossTxNum = append(tctr.totCrossTxNum, 0)

		tctr.broker1TxNum = append(tctr.broker1TxNum, 0)
		tctr.broker2TxNum = append(tctr.broker2TxNum, 0)
		tctr.normalTxNum = append(tctr.normalTxNum, 0)

		tctr.epochID++
	}

	tctr.broker1TxNum[epochid] += b1TxNum
	tctr.broker2TxNum[epochid] += b2TxNum
	tctr.normalTxNum[epochid] += len(b.InnerShardTxs)

	tctr.totCrossTxNum[epochid] += float64(b1TxNum+b2TxNum) / 2
	tctr.totTxNum[epochid] += float64(len(b.InnerShardTxs)) + float64(b1TxNum+b2TxNum)/2
}

func (tctr *TestCrossTxRate_Broker) HandleExtraMessage([]byte) {}

func (tctr *TestCrossTxRate_Broker) OutputRecord() (perEpochCTXratio []float64, totCTXratio float64) {
	tctr.writeToCSV()

	// calculate the simple result
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

func (tctr *TestCrossTxRate_Broker) writeToCSV() {
	fileName := tctr.OutputMetricName()
	measureName := []string{"EpochID", "Total tx # in this epoch", "CTX # in this epoch", "Normal tx # in this epoch", "Broker1 tx # in this epoch", "Broker2 tx # in this epoch", "CTX ratio of this epoch"}
	measureVals := make([][]string, 0)

	for eid, totTxInE := range tctr.totTxNum {
		csvLine := []string{
			strconv.Itoa(eid),
			strconv.FormatFloat(totTxInE, 'f', '8', 64),
			strconv.FormatFloat(tctr.totCrossTxNum[eid], 'f', '8', 64),
			strconv.Itoa(tctr.normalTxNum[eid]),
			strconv.Itoa(tctr.broker1TxNum[eid]),
			strconv.Itoa(tctr.broker2TxNum[eid]),
			strconv.FormatFloat(tctr.totCrossTxNum[eid]/totTxInE, 'f', '8', 64),
		}
		measureVals = append(measureVals, csvLine)
	}
	WriteMetricsToCSV(fileName, measureName, measureVals)
}
