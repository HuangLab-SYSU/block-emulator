package utils

import (
	"blockEmulator/params"
	"log"
	"strconv"
)

// the default method
func Addr2Shard(addr Address) int {
	last8_addr := addr
	if len(last8_addr) > 8 {
		last8_addr = last8_addr[len(last8_addr)-8:]
	}
	num, err := strconv.ParseUint(last8_addr, 16, 64)
	if err != nil {
		log.Panic(err)
	}
	return int(num) % params.ShardNum
}
