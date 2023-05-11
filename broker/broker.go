package broker

import (
	"blockEmulator/message"
	"blockEmulator/params"
	"strconv"
)

type Broker struct {
	BrokerRawMegs  map[string]*message.BrokerRawMeg
	ChainConfig    *params.ChainConfig
	BrokerAddress  []string
	RawTx2BrokerTx map[string][]string
}

func (b *Broker) NewBroker(pcc *params.ChainConfig) {
	b.BrokerRawMegs = make(map[string]*message.BrokerRawMeg)
	b.RawTx2BrokerTx = make(map[string][]string)
	b.ChainConfig = pcc
	b.BrokerAddress = make([]string, 0)
	for i := 0; i < params.ShardNum; i++ {
		tempAddr := "0x000000000000000000000"
		tempAddr = tempAddr + strconv.Itoa(i)
		b.BrokerAddress = append(b.BrokerAddress, tempAddr)
	}
}
