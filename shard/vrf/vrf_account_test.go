package vrf

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	dataDir    = "../test/"
	vrfAccount *VrfAccount
	testmsg    = hexutil.MustDecode("0xce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	sig        []byte
)

func TestVrfAccount(t *testing.T) {
	// Test new
	vrfAccount = NewVrfAccount(dataDir)
	printAccounts(vrfAccount)

	// Test sign
	sig = vrfAccount.SignHash(testmsg)
	if !VerifySignature(testmsg, sig, *vrfAccount.GetAccountAddress()) {
		t.Error("wrong verify sig")
	}

	// Test generate & verify
	seed, err := RlpHash("random seed")
	if err != nil {
		t.Error("rlpHash fail")
	}
	vrfResult := vrfAccount.GenerateVRFOutput(seed[:])
	valid := vrfAccount.VerifyVRFOutput(vrfResult, seed[:])
	if !valid {
		t.Error("verify vrf fail.")
	}
}
