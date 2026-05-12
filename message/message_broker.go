package message

import (
	"blockEmulator/core"
	"blockEmulator/utils"
)

var (
	BrokerRawTx    MessageType = "brokerRawTx"
	BrokerConfirm1 MessageType = "brokerConfirm1"
	BrokerConfirm2 MessageType = "brokerConfirm2"
	BrokerType1    MessageType = "brokerType1"
	BrokerType2    MessageType = "brokerType2"

	CInjectBroker MessageType = "InjectTx_Broker"

	CBrokerTxMap MessageType = "BrokerTxMap"

	CAccountTransferMsg_broker MessageType = "BrokerAS_transfer"
	CInner2CrossTx             MessageType = "innerShardTx_be_crossShard"
)

type BrokerRawMeg struct {
	Tx        *core.Transaction
	Broker    utils.Address
	Hlock     uint64 //ignore
	Snonce    uint64 //ignore
	Bnonce    uint64 //ignore
	Signature []byte // not implemented now.
}

type BrokerType1Meg struct {
	RawMeg   *BrokerRawMeg
	Hcurrent uint64        //ignore
	Broker   utils.Address // replace signature of broker
}

type Mag1Confirm struct {
	//RawDigest string
	Tx1Hash []byte
	RawMeg  *BrokerRawMeg
}

type BrokerType2Meg struct {
	RawMeg *BrokerRawMeg
	Broker utils.Address // replace signature of broker
}

type Mag2Confirm struct {
	//RawDigest string
	Tx2Hash []byte
	RawMeg  *BrokerRawMeg
}

type BrokerTxMap struct {
	BrokerTx2Broker12 map[string][]string // map: raw broker tx to its broker1Tx and broker2Tx
}

type InnerTx2CrossTx struct {
	Txs []*core.Transaction // if an inner-shard tx becomes a cross-shard tx, it will be added into here.
}
