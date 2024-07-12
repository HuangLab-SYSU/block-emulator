package measure

import (
	"blockEmulator/message"
	"strconv"
)

// to test cross-transaction rate
type TestCrossTxRate_Relay struct {
	epochID int

	normalTxNum []int
	relay1TxNum []int
	relay2TxNum []int

	totTxNum      []float64
	totCrossTxNum []float64
}

func NewTestCrossTxRate_Relay() *TestCrossTxRate_Relay {
	return &TestCrossTxRate_Relay{
		epochID:       -1,
		totTxNum:      make([]float64, 0),
		totCrossTxNum: make([]float64, 0),

		normalTxNum: make([]int, 0),
		relay1TxNum: make([]int, 0),
		relay2TxNum: make([]int, 0),
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

		tctr.relay1TxNum = append(tctr.relay1TxNum, 0)
		tctr.relay2TxNum = append(tctr.relay2TxNum, 0)
		tctr.normalTxNum = append(tctr.normalTxNum, 0)

		tctr.epochID++
	}

	tctr.normalTxNum[epochid] += len(b.InnerShardTxs)
	tctr.relay1TxNum[epochid] += r1TxNum
	tctr.relay2TxNum[epochid] += r2TxNum

	tctr.totCrossTxNum[epochid] += float64(r1TxNum+r2TxNum) / 2
	tctr.totTxNum[epochid] += float64(r1TxNum+r2TxNum)/2 + float64(len(b.InnerShardTxs))
}

func (tctr *TestCrossTxRate_Relay) HandleExtraMessage([]byte) {}

func (tctr *TestCrossTxRate_Relay) OutputRecord() (perEpochCTXratio []float64, totCTXratio float64) {
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

func (tctr *TestCrossTxRate_Relay) writeToCSV() {
	fileName := tctr.OutputMetricName()
	measureName := []string{"EpochID", "Total tx # in this epoch", "CTX # in this epoch", "Normal tx # in this epoch", "Relay1 tx # in this epoch", "Relay2 tx # in this epoch", "CTX ratio of this epoch"}
	measureVals := make([][]string, 0)

	for eid, totTxInE := range tctr.totTxNum {
		csvLine := []string{
			strconv.Itoa(eid),
			strconv.FormatFloat(totTxInE, 'f', '8', 64),
			strconv.FormatFloat(tctr.totCrossTxNum[eid], 'f', '8', 64),
			strconv.Itoa(tctr.normalTxNum[eid]),
			strconv.Itoa(tctr.relay1TxNum[eid]),
			strconv.Itoa(tctr.relay2TxNum[eid]),
			strconv.FormatFloat(tctr.totCrossTxNum[eid]/totTxInE, 'f', '8', 64),
		}
		measureVals = append(measureVals, csvLine)
	}
	WriteMetricsToCSV(fileName, measureName, measureVals)
}
