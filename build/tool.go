package build

import (
	"blockEmulator/params"
	"fmt"
	"log"
	"os"
)

func GenerateBatFile(nodenum, shardnum, modID int, dataRootDir string) {
	fileName := fmt.Sprintf("bat_shardNum=%v_NodeNum=%v_mod=%v.bat", shardnum, nodenum, params.CommitteeMethod[modID])
	ofile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Panic(err)
	}
	defer ofile.Close()
	for i := 1; i < nodenum; i++ {
		for j := 0; j < shardnum; j++ {
			str := fmt.Sprintf("start cmd /k go run main.go -n %d -N %d -s %d -S %d -m %d -d %s \n\n", i, nodenum, j, shardnum, modID, dataRootDir)
			ofile.WriteString(str)
		}
	}

	for j := 0; j < shardnum; j++ {
		str := fmt.Sprintf("start cmd /k go run main.go -n 0 -N %d -s %d -S %d -m %d -d %s \n\n", nodenum, j, shardnum, modID, dataRootDir)
		ofile.WriteString(str)
	}

	str := fmt.Sprintf("start cmd /k go run main.go -c -N %d -S %d -m %d -d %s \n\n", nodenum, shardnum, modID, dataRootDir)

	ofile.WriteString(str)
}

func GenerateShellFile(nodenum, shardnum, modID int, dataRootDir string) {
	fileName := fmt.Sprintf("bat_shardNum=%v_NodeNum=%v_mod=%v.sh", shardnum, nodenum, params.CommitteeMethod[modID])
	ofile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Panic(err)
	}
	defer ofile.Close()
	ofile.WriteString("#!/bin/bash \n\n")
	for j := 0; j < shardnum; j++ {
		for i := 1; i < nodenum; i++ {
			str := fmt.Sprintf("go run main.go -n %d -N %d -s %d -S %d -m %d -d %s &\n\n", i, nodenum, j, shardnum, modID, dataRootDir)
			ofile.WriteString(str)
		}
	}

	for j := 0; j < shardnum; j++ {
		str := fmt.Sprintf("go run main.go -n 0 -N %d -s %d -S %d -m %d -d %s &\n\n", nodenum, j, shardnum, modID, dataRootDir)
		ofile.WriteString(str)
	}

	str := fmt.Sprintf("go run main.go -c -N %d -S %d -m %d -d %s &\n\n", nodenum, shardnum, modID, dataRootDir)

	ofile.WriteString(str)
}

func Exebat_Windows_GenerateBatFile(nodenum, shardnum, modID int, dataRootDir string) {
	fileName := fmt.Sprintf("WinExe_bat_shardNum=%v_NodeNum=%v_mod=%v.bat", shardnum, nodenum, params.CommitteeMethod[modID])
	ofile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Panic(err)
	}
	defer ofile.Close()
	for i := 1; i < nodenum; i++ {
		for j := 0; j < shardnum; j++ {
			str := fmt.Sprintf("start cmd /k blockEmulator_Windows_Precompile.exe -n %d -N %d -s %d -S %d -m %d -d %s\n\n", i, nodenum, j, shardnum, modID, dataRootDir)
			ofile.WriteString(str)
		}
	}

	for j := 0; j < shardnum; j++ {
		str := fmt.Sprintf("start cmd /k blockEmulator_Windows_Precompile.exe -n 0 -N %d -s %d -S %d -m %d -d %s\n\n", nodenum, j, shardnum, modID, dataRootDir)
		ofile.WriteString(str)
	}

	str := fmt.Sprintf("start cmd /k blockEmulator_Windows_Precompile.exe -c -N %d -S %d -m %d -d %s\n\n", nodenum, shardnum, modID, dataRootDir)

	ofile.WriteString(str)
}

func Exebat_Linux_GenerateShellFile(nodenum, shardnum, modID int, dataRootDir string) {
	fileName := fmt.Sprintf("bat_shardNum=%v_NodeNum=%v_mod=%v.sh", shardnum, nodenum, params.CommitteeMethod[modID])
	ofile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Panic(err)
	}
	defer ofile.Close()
	ofile.WriteString("#!/bin/bash \n\n")
	for j := 0; j < shardnum; j++ {
		for i := 1; i < nodenum; i++ {
			str := fmt.Sprintf("./blockEmulator_Linux_Precompile -n %d -N %d -s %d -S %d -m %d -d %s &\n\n", i, nodenum, j, shardnum, modID, dataRootDir)
			ofile.WriteString(str)
		}
	}

	for j := 0; j < shardnum; j++ {
		str := fmt.Sprintf("./blockEmulator_Linux_Precompile -n 0 -N %d -s %d -S %d -m %d -d %s &\n\n", nodenum, j, shardnum, modID, dataRootDir)
		ofile.WriteString(str)
	}

	str := fmt.Sprintf("./blockEmulator_Linux_Precompile -c -N %d -S %d -m %d -d %s &\n\n", nodenum, shardnum, modID, dataRootDir)

	ofile.WriteString(str)
}

func Exebat_MacOS_GenerateShellFile(nodenum, shardnum, modID int, dataRootDir string) {
	fileName := fmt.Sprintf("bat_shardNum=%v_NodeNum=%v_mod=%v.sh", shardnum, nodenum, params.CommitteeMethod[modID])
	ofile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Panic(err)
	}
	defer ofile.Close()
	ofile.WriteString("#!/bin/bash \n\n")
	for j := 0; j < shardnum; j++ {
		for i := 1; i < nodenum; i++ {
			str := fmt.Sprintf("./blockEmulator_MacOS_Precompile -n %d -N %d -s %d -S %d -m %d -d %s &\n\n", i, nodenum, j, shardnum, modID, dataRootDir)
			ofile.WriteString(str)
		}
	}

	for j := 0; j < shardnum; j++ {
		str := fmt.Sprintf("./blockEmulator_MacOS_Precompile -n 0 -N %d -s %d -S %d -m %d -d %s &\n\n", nodenum, j, shardnum, modID, dataRootDir)
		ofile.WriteString(str)
	}

	str := fmt.Sprintf("./blockEmulator_MacOS_Precompile -c -N %d -S %d -m %d -d %s &\n\n", nodenum, shardnum, modID, dataRootDir)

	ofile.WriteString(str)
}
