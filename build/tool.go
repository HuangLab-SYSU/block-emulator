package build

import (
	"blockEmulator/params"
	"fmt"
	"log"
	"os"
)

func GenerateBatFile(nodenum, shardnum, modID int) {
	fileName := fmt.Sprintf("bat_shardNum=%v_NodeNum=%v_mod=%v.bat", shardnum, nodenum, params.CommitteeMethod[modID])
	ofile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Panic(err)
	}
	defer ofile.Close()
	for i := 1; i < nodenum; i++ {
		for j := 0; j < shardnum; j++ {
			str := fmt.Sprintf("start cmd /k go run main.go -n %d -N %d -s %d -S %d -m %d \n\n", i, nodenum, j, shardnum, modID)
			ofile.WriteString(str)
		}
	}

	for j := 0; j < shardnum; j++ {
		str := fmt.Sprintf("start cmd /k go run main.go -n 0 -N %d -s %d -S %d -m %d \n\n", nodenum, j, shardnum, modID)
		ofile.WriteString(str)
	}

	str := fmt.Sprintf("start cmd /k go run main.go -c -N %d -S %d -m %d \n\n", nodenum, shardnum, modID)

	ofile.WriteString(str)
}

func GenerateShellFile(nodenum, shardnum, modID int) {
	fileName := fmt.Sprintf("bat_shardNum=%v_NodeNum=%v_mod=%v.sh", shardnum, nodenum, params.CommitteeMethod[modID])
	ofile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Panic(err)
	}
	defer ofile.Close()
	str1 := fmt.Sprintf("#!/bin/bash \n\n")
	ofile.WriteString(str1)
	for j := 0; j < shardnum; j++ {
		for i := 1; i < nodenum; i++ {
			str := fmt.Sprintf("go run main.go -n %d -N %d -s %d -S %d -m %d &\n\n", i, nodenum, j, shardnum, modID)
			ofile.WriteString(str)
		}
	}

	for j := 0; j < shardnum; j++ {
		str := fmt.Sprintf("go run main.go -n 0 -N %d -s %d -S %d -m %d &\n\n", nodenum, j, shardnum, modID)
		ofile.WriteString(str)
	}

	str := fmt.Sprintf("go run main.go -c -N %d -S %d -m %d &\n\n", nodenum, shardnum, modID)

	ofile.WriteString(str)
}
