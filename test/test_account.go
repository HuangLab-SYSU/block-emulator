package test

import (
	"blockEmulator/core"
	"fmt"
)

func Test_account() {
	for i := 0; i < 10; i++ {
		address := core.GenerateAddress()
		fmt.Printf("%v\n", address)
	}
}
