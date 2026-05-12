package Broker2Earn

import (
	"blockEmulator/core"
	"blockEmulator/message"
	"blockEmulator/params"
	"blockEmulator/utils"
	"math/big"
)

type RatioBrokerRawMeg struct {
	Tx          *core.Transaction
	BrokerRatio map[utils.Address]float64
}

// B2E rounding core functin
// Input tx list and brokerTable (broker balance in each shard)
// Output a map that ctx -> broker
func B2E(brokerRawMegs []*message.BrokerRawMeg, BrokerBalance map[string]map[uint64]*big.Int) ([]*message.BrokerRawMeg, []*message.BrokerRawMeg) {
	// 1. find ctx
	New_brokerRawMegs := make([][]*message.BrokerRawMeg, 2)
	New_brokerRawMegs[0] = make([]*message.BrokerRawMeg, 0)
	New_brokerRawMegs[1] = make([]*message.BrokerRawMeg, 0)
	sign := 0
	for _, brokerRawMeg := range brokerRawMegs {
		if utils.Addr2Shard(brokerRawMeg.Tx.Recipient) != utils.Addr2Shard(brokerRawMeg.Tx.Sender) {
			New_brokerRawMegs[sign] = append(New_brokerRawMegs[sign], brokerRawMeg)
		}
	}

	// 2.accumulate broker balance

	brokerBalance := make(map[string]*big.Int)
	nowbrokerBalance := make(map[string]*big.Int)

	for brokerID, shardAmounts := range BrokerBalance {
		totalAmount := big.NewInt(0)
		for _, amount := range shardAmounts {
			totalAmount.Add(totalAmount, amount)
		}
		brokerBalance[brokerID] = totalAmount
		nowbrokerBalance[brokerID] = new(big.Int).Set(brokerBalance[brokerID])
	}

	result := make([]*message.BrokerRawMeg, 0)

	// Run iter
	for iter := 0; iter < params.IterNum_B2E; iter++ {
		//println("New_brokerRawMegs length is ", len(New_brokerRawMegs[sign]))
		// Run b2e rounding
		// run URFA_linear
		URFA_Linear_result := URFA_Linear(New_brokerRawMegs[sign], nowbrokerBalance)
		rounding_result := B2E_Rounding(URFA_Linear_result, nowbrokerBalance)

		//println("rounding result ", len(rounding_result))
		for _, BrokerRawMeg := range rounding_result {
			//println("Allocated tx ", index, " ", BrokerRawMeg.Broker, " ", BrokerRawMeg.Tx.Value.String(), " ", BrokerRawMeg.Tx.Fee.String(), " ", BrokerRawMeg.Tx.TxHash, " ", BrokerRawMeg.Tx.Sender, " ", BrokerRawMeg.Tx.Recipient)
			nowbrokerBalance[BrokerRawMeg.Broker].Sub(nowbrokerBalance[BrokerRawMeg.Broker], BrokerRawMeg.Tx.Value)
			//println("Iter ", iter, " ", BrokerRawMeg.Broker, " ", nowbrokerBalance[BrokerRawMeg.Broker].String())
		}
		result = append(result, rounding_result...)

		New_brokerRawMegs[1-sign] = difference(New_brokerRawMegs[sign], rounding_result)
		sign = 1 - sign
	}
	//for brokeraddress, balance := range nowbrokerBalance {
	//	println("result ", brokeraddress, " ", balance.String())
	//}
	restBrokerRawMeg := difference(brokerRawMegs, result)

	return result, restBrokerRawMeg
}
func difference(slice1, slice2 []*message.BrokerRawMeg) []*message.BrokerRawMeg {
	m := make(map[*core.Transaction]bool)
	for _, item := range slice2 {
		m[item.Tx] = true
	}
	var diff []*message.BrokerRawMeg
	for _, item := range slice1 {
		if !m[item.Tx] {
			diff = append(diff, item)
		}
	}
	return diff
}
