package Broker2Earn

import (
	"blockEmulator/core"
	"blockEmulator/message"
	"blockEmulator/utils"
	"math/big"
	"math/rand"
	"time"
)

// Broker2earn rounding functin
// Input brokerTable (brokers balance) and tx allcolated result by relax function
// Output a map that ctx -> broker

func B2E_Rounding(RatioBrokerRawMegs []*RatioBrokerRawMeg, BrokerBalance map[string]*big.Int) []*message.BrokerRawMeg {

	result := make([]*message.BrokerRawMeg, 0)
	// build tmp broker balance
	nowBrokerBalance := make(map[string]*big.Int)
	for brokerID, _ := range BrokerBalance {
		nowBrokerBalance[brokerID] = new(big.Int).SetInt64(0)
	}
	// get all broker address
	brokerAddresses := make([]string, 0, len(BrokerBalance))

	for brokerAddress := range BrokerBalance {
		brokerAddresses = append(brokerAddresses, brokerAddress)
	}

	// 1. build map broker address -> tx & ratio, tx -? broker
	broker2Tx := make(map[utils.Address][]*core.Transaction)
	broker2Ratio := make(map[utils.Address][]float64)
	tx2broker := make(map[*core.Transaction]map[utils.Address]int)

	for _, brokerAddress := range brokerAddresses {
		broker2Tx[brokerAddress] = make([]*core.Transaction, 0)
		broker2Ratio[brokerAddress] = make([]float64, 0)
	}

	for _, ratioBrokerRawMeg := range RatioBrokerRawMegs {
		for brokerAddress, ratio := range ratioBrokerRawMeg.BrokerRatio {
			broker2Tx[brokerAddress] = append(broker2Tx[brokerAddress], ratioBrokerRawMeg.Tx)
			broker2Ratio[brokerAddress] = append(broker2Ratio[brokerAddress], ratio)
			if tx2broker[ratioBrokerRawMeg.Tx] == nil {
				tx2broker[ratioBrokerRawMeg.Tx] = make(map[utils.Address]int)
			}
			tx2broker[ratioBrokerRawMeg.Tx][brokerAddress] = len(broker2Tx[brokerAddress]) - 1
		}
	}
	// 2. selected tx
	for _, brokerAddress := range brokerAddresses {
		for len(broker2Tx[brokerAddress]) != 0 {
			p := broker2Ratio[brokerAddress]
			var sumP float64
			for _, v := range p {
				sumP += v
			}
			sumP += 0.0000000000000000001
			for i := range p {
				p[i] /= sumP
			}
			choosed_tx_id := 0

			if len(broker2Tx[brokerAddress]) > 1 {
				choosed_tx_id = Random_choice(p)
			}
			tmpValue := new(big.Int).SetInt64(0)
			tmpValue.Add(nowBrokerBalance[brokerAddress], broker2Tx[brokerAddress][choosed_tx_id].Value)
			//println("IterX ", brokerAddress, " ", broker2Tx[brokerAddress][choosed_tx_id].Value.String(), " ", tmpValue.String(), " ", BrokerBalance[brokerAddress].String())
			if tmpValue.Cmp(BrokerBalance[brokerAddress]) <= 0 {
				tmp_brokerRawMeg := &message.BrokerRawMeg{
					Tx:     broker2Tx[brokerAddress][choosed_tx_id],
					Broker: brokerAddress,
				}
				nowBrokerBalance[brokerAddress].Add(nowBrokerBalance[brokerAddress], broker2Tx[brokerAddress][choosed_tx_id].Value)
				result = append(result, tmp_brokerRawMeg)
			}
			choosed_tx := broker2Tx[brokerAddress][choosed_tx_id]
			indexTx2Broker := tx2broker[choosed_tx]
			for bd, index := range indexTx2Broker {
				broker2Tx[bd] = append(broker2Tx[bd][:index], broker2Tx[bd][index+1:]...)
				broker2Ratio[bd] = append(broker2Ratio[bd][:index], broker2Ratio[bd][index+1:]...)
				for _, tx := range broker2Tx[bd][index:] {
					tx2broker[tx][bd] -= 1
					//println("Tx ", tx.TxHash, " ", bd, " ", tx2broker[tx][bd])
				}
			}

		}
	}
	return result
}

func Random_choice(probility []float64) int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Float64()

	// 初始化累积概率
	cumulativeProb := 0.0

	// 遍历概率分布数组，根据随机数选择相应的值
	for i, prob := range probility {
		cumulativeProb += prob
		if r <= cumulativeProb {
			return i
		}
	}
	return -1
}
