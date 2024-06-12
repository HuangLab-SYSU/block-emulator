package test

import (
	"blockEmulator/account"
	"fmt"
)

func Test_account() {
	for i := 0; i < 10; i++ {
		address := account.GenerateAddress()
		fmt.Printf("%v\n", address)
	}
}
