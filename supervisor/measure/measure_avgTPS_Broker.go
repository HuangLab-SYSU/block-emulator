package measure

import (
	"blockEmulator/message"
	"time"
)

type TestModule_avgTPS_Broker struct {
	epochID      int
	excutedTxNum []float64   // record how many excuted txs in a epoch, maybe the cross shard tx will be calculated as a 0.5 tx
	startTime    []time.Time // record when the epoch starts
	endTime      []time.Time // record when the epoch ends
}

func NewTestModule_avgTPS_Broker() *TestModule_avgTPS_Broker {
	return &TestModule_avgTPS_Broker{
		epochID:      -1,
		excutedTxNum: make([]float64, 0),
		startTime:    make([]time.Time, 0),
		endTime:      make([]time.Time, 0),
	}
}

func (tat *TestModule_avgTPS_Broker) OutputMetricName() string {
	return "Average_TPS"
}

func (tat *TestModule_avgTPS_Broker) OutputMetricTitle() string {
	return "Average_TPS"
}

// add the number of excuted txs, and change the time records
func (tat *TestModule_avgTPS_Broker) UpdateMeasureRecord(b *message.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}

	epochid := b.Epoch
	txNum := float64(len(b.ExcutedTxs))
	earliestTime := b.ProposeTime
	latestTime := b.CommitTime

	// extend
	for tat.epochID < epochid {
		tat.excutedTxNum = append(tat.excutedTxNum, 0)
		tat.startTime = append(tat.startTime, time.Time{})
		tat.endTime = append(tat.endTime, time.Time{})
		tat.epochID++
	}
	// modify the local epoch
	tat.excutedTxNum[epochid] += txNum + (float64(b.Broker1TxNum)+float64(b.Broker2TxNum))/2
	if tat.startTime[epochid].IsZero() || tat.startTime[epochid].After(earliestTime) {
		tat.startTime[epochid] = earliestTime
	}
	if tat.endTime[epochid].IsZero() || latestTime.After(tat.endTime[epochid]) {
		tat.endTime[epochid] = latestTime
	}
}

func (tat *TestModule_avgTPS_Broker) HandleExtraMessage([]byte) {}

// output the average TPS
func (tat *TestModule_avgTPS_Broker) OutputRecord() (perEpochTPS []float64, totalTPS float64) {
	perEpochTPS = make([]float64, tat.epochID+1)
	totalTxNum := 0.0
	eTime := time.Now()
	lTime := time.Time{}
	for eid, exTxNum := range tat.excutedTxNum {
		timeGap := tat.endTime[eid].Sub(tat.startTime[eid]).Seconds()
		perEpochTPS[eid] = exTxNum / timeGap
		totalTxNum += exTxNum
		if eTime.After(tat.startTime[eid]) {
			eTime = tat.startTime[eid]
		}
		if tat.endTime[eid].After(lTime) {
			lTime = tat.endTime[eid]
		}
	}
	totalTPS = totalTxNum / (lTime.Sub(eTime).Seconds())
	return
}
