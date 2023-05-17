package broker

import (
	"blockEmulator/message"
	"blockEmulator/params"
	"bufio"
	"fmt"
	"os"
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
	b.BrokerAddress = b.initBrokerAddr(params.BrokerNum)
}

func (b *Broker) IsBroker(address string) bool {
	for _, brokerAddress := range b.BrokerAddress {
		if brokerAddress == address {
			return true
		}
	}
	return false
}

func (b *Broker) initBrokerAddr(num int) []string {
	brokerAddress := make([]string, 0)
	filePath := `./broker/broker`
	readFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		brokerAddress = append(brokerAddress, fileScanner.Text())
		num--
		if num == 0 {
			break
		}
	}
	readFile.Close()
	return brokerAddress
}
