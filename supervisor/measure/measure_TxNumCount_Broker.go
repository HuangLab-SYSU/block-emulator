package measure

import (
	"blockEmulator/message"
	"strconv"
)

// to test Tx number
type TestTxNumCount_Broker struct {
	epochID      int
	txNum        []float64
	broker1TxNum []int // record how many broker1 txs in an epoch.
	broker2TxNum []int // record how many broker2 txs in an epoch.
	normalTxNum  []int // record how many normal txs in an epoch.
}

func NewTestTxNumCount_Broker() *TestTxNumCount_Broker {
	return &TestTxNumCount_Broker{
		epochID: -1,
		txNum:   make([]float64, 0),

		broker1TxNum: make([]int, 0),
		broker2TxNum: make([]int, 0),
		normalTxNum:  make([]int, 0),
	}
}

func (ttnc *TestTxNumCount_Broker) OutputMetricName() string {
	return "Tx_number"
}

func (ttnc *TestTxNumCount_Broker) UpdateMeasureRecord(b *message.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}
	epochid := b.Epoch
	b1TxNum := len(b.Broker1Txs)
	b2TxNum := len(b.Broker2Txs)
	// extend
	for ttnc.epochID < epochid {
		ttnc.txNum = append(ttnc.txNum, 0)
		ttnc.broker1TxNum = append(ttnc.broker1TxNum, 0)
		ttnc.broker2TxNum = append(ttnc.broker2TxNum, 0)
		ttnc.normalTxNum = append(ttnc.normalTxNum, 0)

		ttnc.epochID++
	}

	ttnc.broker1TxNum[epochid] += b1TxNum
	ttnc.broker2TxNum[epochid] += b2TxNum
	ttnc.normalTxNum[epochid] += len(b.InnerShardTxs)
	ttnc.txNum[epochid] += float64(len(b.InnerShardTxs)) + (float64(b1TxNum)+float64(b2TxNum))/2
}

func (ttnc *TestTxNumCount_Broker) HandleExtraMessage([]byte) {}

func (ttnc *TestTxNumCount_Broker) OutputRecord() (perEpochCTXs []float64, totTxNum float64) {
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

func (ttnc *TestTxNumCount_Broker) writeToCSV() {
	fileName := ttnc.OutputMetricName()
	measureName := []string{"EpochID", "Total tx # in this epoch", "Normal tx # in this epoch", "Broker1 tx # in this epoch", "Broker2 tx # in this epoch"}
	measureVals := make([][]string, 0)

	for eid, totTxInE := range ttnc.txNum {
		csvLine := []string{
			strconv.Itoa(eid),
			strconv.FormatFloat(totTxInE, 'f', '8', 64),
			strconv.Itoa(ttnc.normalTxNum[eid]),
			strconv.Itoa(ttnc.broker1TxNum[eid]),
			strconv.Itoa(ttnc.broker2TxNum[eid]),
		}
		measureVals = append(measureVals, csvLine)
	}
	WriteMetricsToCSV(fileName, measureName, measureVals)
}
