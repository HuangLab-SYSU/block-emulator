package measure

import (
	"blockEmulator/message"
	"strconv"
)

// to test Tx number
type TestTxNumCount_Relay struct {
	epochID int
	txNum   []float64

	normalTxNum []int
	relay1TxNum []int
	relay2TxNum []int
}

func NewTestTxNumCount_Relay() *TestTxNumCount_Relay {
	return &TestTxNumCount_Relay{
		epochID: -1,
		txNum:   make([]float64, 0),

		normalTxNum: make([]int, 0),
		relay1TxNum: make([]int, 0),
		relay2TxNum: make([]int, 0),
	}
}

func (ttnc *TestTxNumCount_Relay) OutputMetricName() string {
	return "Tx_number"
}

func (ttnc *TestTxNumCount_Relay) UpdateMeasureRecord(b *message.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}
	epochid := b.Epoch
	r1TxNum := len(b.Relay1Txs)
	r2TxNum := len(b.Relay2Txs)
	// extend
	for ttnc.epochID < epochid {
		ttnc.txNum = append(ttnc.txNum, 0)
		ttnc.relay1TxNum = append(ttnc.relay1TxNum, 0)
		ttnc.relay2TxNum = append(ttnc.relay2TxNum, 0)
		ttnc.normalTxNum = append(ttnc.normalTxNum, 0)

		ttnc.epochID++
	}

	ttnc.normalTxNum[epochid] += len(b.InnerShardTxs)
	ttnc.relay1TxNum[epochid] += r1TxNum
	ttnc.relay2TxNum[epochid] += r2TxNum
	ttnc.txNum[epochid] += float64(len(b.InnerShardTxs)) + float64(len(b.Relay1Txs)+len(b.Relay2Txs))/2
}

func (ttnc *TestTxNumCount_Relay) HandleExtraMessage([]byte) {}

func (ttnc *TestTxNumCount_Relay) OutputRecord() (perEpochCTXs []float64, totTxNum float64) {
	ttnc.writeToCSV()

	// calculate the simple result
	perEpochCTXs = make([]float64, 0)
	totTxNum = 0.0
	for _, tn := range ttnc.txNum {
		perEpochCTXs = append(perEpochCTXs, tn)
		totTxNum += tn
	}
	return perEpochCTXs, totTxNum
}

func (ttnc *TestTxNumCount_Relay) writeToCSV() {
	fileName := ttnc.OutputMetricName()
	measureName := []string{"EpochID", "Total tx # in this epoch", "Normal tx # in this epoch", "Relay1 tx # in this epoch", "Relay2 tx # in this epoch"}
	measureVals := make([][]string, 0)

	for eid, totTxInE := range ttnc.txNum {
		csvLine := []string{
			strconv.Itoa(eid),
			strconv.FormatFloat(totTxInE, 'f', '8', 64),
			strconv.Itoa(ttnc.normalTxNum[eid]),
			strconv.Itoa(ttnc.relay1TxNum[eid]),
			strconv.Itoa(ttnc.relay2TxNum[eid]),
		}
		measureVals = append(measureVals, csvLine)
	}
	WriteMetricsToCSV(fileName, measureName, measureVals)
}
