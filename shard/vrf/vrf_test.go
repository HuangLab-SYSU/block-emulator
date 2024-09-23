package vrf

import (
	"fmt"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/log"
)

func TestVRF(t *testing.T) {
	// 选择椭圆曲线，这里选择 secp256k1 曲线
	s, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		log.Error("generate private key fail", "err", err)
	}
	privateKey := s.ToECDSA()

	for i := 0; i < 3; i++ {
		// 构造输入数据
		inputData := []byte(fmt.Sprintf("This is some input data %d", i))

		// 进行 VRF 计算
		vrfResult := GenerateVRF(privateKey, inputData)

		// 输出 VRF 结果
		fmt.Printf("VRF Proof: %x\n", vrfResult.Proof)
		fmt.Printf("Random Value: %x\n", vrfResult.RandomValue)

		// 验证 VRF 结果
		isValid := VerifyVRF(&privateKey.PublicKey, inputData, vrfResult)
		fmt.Println("VRF Verification:", isValid)
	}
}
