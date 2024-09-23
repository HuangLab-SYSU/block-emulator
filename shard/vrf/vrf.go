package vrf

import (
	"crypto/ecdsa"
	"reflect"

	"github.com/ethereum/go-ethereum/log"
	"github.com/vechain/go-ecvrf"
)

// VRFResult 包含 VRF 方法的输出结果
type VRFResult struct {
	Proof       []byte // VRF 证明
	RandomValue []byte // 随机数
}

// GenerateVRF 使用私钥进行 VRF 计算
//
// VRF算法：采用第三方库ecvrf
func GenerateVRF(privateKey *ecdsa.PrivateKey, input []byte) *VRFResult {
	output, proof, err := ecvrf.Secp256k1Sha256Tai.Prove(privateKey, input)
	if err != nil {
		log.Error("GenerateVRF fail", "err", err)
	}
	return &VRFResult{
		Proof:       proof,
		RandomValue: output,
	}
}

// VerifyVRF 使用公钥进行 VRF 结果验证
//
// VRF算法：采用第三方库ecvrf
func VerifyVRF(publicKey *ecdsa.PublicKey, input []byte, vrfResult *VRFResult) bool {
	output, err := ecvrf.Secp256k1Sha256Tai.Verify(publicKey, input, vrfResult.Proof)
	if err != nil {
		log.Error("VerifyVRF fail", "err", err)
	}

	return reflect.DeepEqual(output, vrfResult.RandomValue)
}
