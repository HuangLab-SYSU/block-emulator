package build

import (
	"blockEmulator/params"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
)

var absolute_path = getABpath()

func getABpath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

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
	str := fmt.Sprintf("start cmd /k go run main.go -c -N %d -S %d -m %d \n\n", nodenum, shardnum, modID)

	ofile.WriteString(str)
	for j := 0; j < shardnum; j++ {
		str := fmt.Sprintf("start cmd /k go run main.go -n 0 -N %d -s %d -S %d -m %d \n\n", nodenum, j, shardnum, modID)
		ofile.WriteString(str)
	}
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
	str := fmt.Sprintf("go run main.go -c -N %d -S %d -m %d &\n\n", nodenum, shardnum, modID)

	ofile.WriteString(str)
	for j := 0; j < shardnum; j++ {
		str := fmt.Sprintf("go run main.go -n 0 -N %d -s %d -S %d -m %d &\n\n", nodenum, j, shardnum, modID)
		ofile.WriteString(str)
	}
}

func GenerateVBSFile() {
	ofile, err := os.OpenFile("batrun_HideWorker.vbs", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Panic(err)
	}
	defer ofile.Close()
	nodenum := params.NodesInShard
	shardnum := params.ShardNum

	ofile.WriteString("Dim shell\nSet shell = CreateObject(\"WScript.Shell\")\n")

	for i := 1; i < nodenum; i++ {
		for j := 0; j < shardnum; j++ {
			str := fmt.Sprintf("shell.Run "+"\"go run "+absolute_path+"/main.go %d %d %d %d\", 0, false \n", i, nodenum, j, shardnum)
			ofile.WriteString(str)
		}
	}
	for j := 0; j < shardnum; j++ {
		str := fmt.Sprintf("shell.Run "+"\"go run "+absolute_path+"/main.go 0 %d %d %d\", 0, false \n", nodenum, j, shardnum)
		ofile.WriteString(str)
	}
	bfile, err := os.OpenFile("batrun_showSupervisor.bat", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Panic(err)
	}
	defer bfile.Close()
	bfile.WriteString("start " + absolute_path + "/batrun_HideWorker.vbs \n")
	str := fmt.Sprintf("start cmd /k \"cd "+absolute_path+" && go run main.go 12345678 %d 0 %d \" \n", nodenum, shardnum)
	bfile.WriteString(str)
}
