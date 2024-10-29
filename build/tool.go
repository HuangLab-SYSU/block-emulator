package build

import (
	"fmt"
	"runtime"
	"strings"
)

func GenerateBatchByIpTable(nodenum, shardnum int) error {
	// read IP table file first
	ipMap := readIpTable("./ipTable.json")

	// determine the formats of commands and fileNames, according to operating system
	var fileNameFormat, commandFormat string
	os := runtime.GOOS
	switch os {
	case "windows":
		fileNameFormat = "complie_run_IpAddr=%s.bat"
		commandFormat = "start cmd /k go run main.go"
	default:
		fileNameFormat = "complie_run_IpAddr=%s.sh"
		commandFormat = "go run main.go"
	}

	// generate file for each ip
	for i := 0; i < shardnum; i++ {
		// if this shard is not existed, return
		if _, shard_exist := ipMap[uint64(i)]; !shard_exist {
			return fmt.Errorf("the shard (shardID = %d) is not existed in the IP Table file", i)
		}
		// if this shard is existed.
		for j := 0; j < nodenum; j++ {
			if nodeIp, node_exist := ipMap[uint64(i)][uint64(j)]; node_exist {
				// attach this command to this file
				ipAddr := strings.Split(nodeIp, ":")[0]
				batFilePath := fmt.Sprintf(fileNameFormat, strings.ReplaceAll(ipAddr, ".", "_"))
				command := fmt.Sprintf(commandFormat+" -n %d -N %d -s %d -S %d\n", j, nodenum, i, shardnum)
				if err := attachLineToFile(batFilePath, command); nil != err {
					return err
				}
			} else {
				return fmt.Errorf("the node (shardID = %d, nodeID = %d) is not existed in the IP Table file", i, j)
			}
		}
	}

	// generate command for supervisor
	if supervisorShard, shard_exist := ipMap[2147483647]; shard_exist {
		if nodeIp, node_exist := supervisorShard[0]; node_exist {
			ipAddr := strings.Split(nodeIp, ":")[0]
			batFilePath := fmt.Sprintf(fileNameFormat, strings.ReplaceAll(ipAddr, ".", "_"))
			supervisor_command := fmt.Sprintf(commandFormat+" -c -N %d -S %d\n", nodenum, shardnum)
			if err := attachLineToFile(batFilePath, supervisor_command); nil != err {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("the supervisor (shardID = 2147483647, nodeID = 0) is not existed in the IP Table file")
}

func GenerateExeBatchByIpTable(nodenum, shardnum int) error {
	// read IP table file first
	ipMap := readIpTable("./ipTable.json")

	// determine the formats of commands and fileNames, according to operating system
	var fileNameFormat, commandFormat string
	os := runtime.GOOS
	switch os {
	case "windows":
		fileNameFormat = os + "_exe_run_IpAddr=%s.bat"
		commandFormat = "start cmd /k blockEmulator_Windows_Precompile.exe"
	default:
		fileNameFormat = os + "_exe_run_IpAddr=%s.sh"
		commandFormat = "./blockEmulator_" + os + "_Precompile"
	}

	// generate file for each ip
	for i := 0; i < shardnum; i++ {
		// if this shard is not existed, return
		if _, shard_exist := ipMap[uint64(i)]; !shard_exist {
			return fmt.Errorf("the shard (shardID = %d) is not existed in the IP Table file", i)
		}
		// if this shard is existed.
		for j := 0; j < nodenum; j++ {
			if nodeIp, node_exist := ipMap[uint64(i)][uint64(j)]; node_exist {
				// attach this command to this file
				ipAddr := strings.Split(nodeIp, ":")[0]
				batFilePath := fmt.Sprintf(fileNameFormat, strings.ReplaceAll(ipAddr, ".", "_"))
				command := fmt.Sprintf(commandFormat+" -n %d -N %d -s %d -S %d\n", j, nodenum, i, shardnum)
				if err := attachLineToFile(batFilePath, command); nil != err {
					return err
				}
			} else {
				return fmt.Errorf("the node (shardID = %d, nodeID = %d) is not existed in the IP Table file", i, j)
			}
		}
	}

	// generate command for supervisor
	if supervisorShard, shard_exist := ipMap[2147483647]; shard_exist {
		if nodeIp, node_exist := supervisorShard[0]; node_exist {
			ipAddr := strings.Split(nodeIp, ":")[0]
			batFilePath := fmt.Sprintf(fileNameFormat, strings.ReplaceAll(ipAddr, ".", "_"))
			supervisor_command := fmt.Sprintf(commandFormat+" -c -N %d -S %d\n", nodenum, shardnum)
			if err := attachLineToFile(batFilePath, supervisor_command); nil != err {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("the supervisor (shardID = 2147483647, nodeID = 0) is not existed in the IP Table file")
}
