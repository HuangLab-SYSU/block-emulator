package vrf

import (
	"crypto/ecdsa"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

type VrfAccount struct {
	privateKey  *ecdsa.PrivateKey
	pubKey      *ecdsa.PublicKey
	accountAddr common.Address
	keyDir      string
}

func NewVrfAccount(nodeDatadir string) *VrfAccount {
	_privateKey := newPrivateKey()

	vrfAccount := &VrfAccount{
		privateKey: _privateKey,
		keyDir:     filepath.Join(nodeDatadir, "KeyStoreDir"),
	}
	vrfAccount.pubKey = &vrfAccount.privateKey.PublicKey
	vrfAccount.accountAddr = crypto.PubkeyToAddress(*vrfAccount.pubKey)
	// fmt.Printf("create addr: %v\n", vrfAccount.accountAddr)
	return vrfAccount
}

func newPrivateKey() *ecdsa.PrivateKey {
	// 选择椭圆曲线，这里选择 secp256k1 曲线
	s, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		log.Error("generate private key fail", "err", err)
	}
	privateKey := s.ToECDSA()
	return privateKey
}

func (vrfAccount *VrfAccount) SignHash(hash []byte) []byte {
	sig, err := crypto.Sign(hash, vrfAccount.privateKey)
	if err != nil {
		log.Error("signHashFail", "err", err)
		return []byte{}
	}
	return sig
}

/*
这个方法是被该账户以外的其他账户调用，以验证签名的正确性的
所以不能直接获取公钥和地址，要从签名中恢复
*/
func VerifySignature(msgHash []byte, sig []byte, expected_addr common.Address) bool {
	// 恢复公钥
	pubKeyBytes, err := crypto.Ecrecover(msgHash, sig)
	if err != nil {
		log.Error("ecrecover fail", "err", err)
		// fmt.Printf("ecrecover err: %v\n", err)
	}

	pubkey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		log.Error("UnmarshalPubkey fail", "err", err)
		// fmt.Printf("UnmarshalPubkey err: %v\n", err)
	}

	recovered_addr := crypto.PubkeyToAddress(*pubkey)
	return recovered_addr == expected_addr
}

/* 接收一个随机种子，用私钥生成一个随机数输出和对应的证明 */
func (vrfAccount *VrfAccount) GenerateVRFOutput(randSeed []byte) *VRFResult {
	// vrfResult := utils.GenerateVRF(vrfAccount.privateKey, randSeed)
	sig := vrfAccount.SignHash(randSeed)
	vrfResult := &VRFResult{
		RandomValue: sig,
		Proof:       vrfAccount.accountAddr[:],
	}
	return vrfResult
}

/* 接收随机数输出和对应证明，用公钥验证该随机数输出是否合法 */
func (vrfAccount *VrfAccount) VerifyVRFOutput(vrfResult *VRFResult, randSeed []byte) bool {
	// return utils.VerifyVRF(vrfAccount.pubKey, randSeed, vrfResult)
	return VerifySignature(randSeed, vrfResult.RandomValue, common.BytesToAddress(vrfResult.Proof))
}

func printAccounts(vrfAccount *VrfAccount) {
	fmt.Println("node account", "keyDir", vrfAccount.keyDir, "address", vrfAccount.accountAddr)
	fmt.Printf("privateKey: %x", crypto.FromECDSA(vrfAccount.privateKey))
}

func (vrfAccount *VrfAccount) GetAccountAddress() *common.Address {
	return &vrfAccount.accountAddr
}

/** 先将结构体进行rlp编码，再哈希
 * 注意: !!!结构体中不能存在int类型的变量！！！
 */
func rlpHash(x interface{}) (h common.Hash, err error) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	err = rlp.Encode(sha, x)
	if err != nil {
		return common.Hash{}, err
	}
	_, err = sha.Read(h[:])
	return h, err
}

func RlpHash(x interface{}) (h common.Hash, err error) {
	return rlpHash(x)
}
