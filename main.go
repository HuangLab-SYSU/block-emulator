package main

import (
	"blockEmulator/build"
	"blockEmulator/params"
	"runtime"

	"github.com/spf13/pflag"
)

var (
	// network config
	shardNum int
	nodeNum  int
	shardID  int
	nodeID   int

	// consensus module
	modID int

	// the location for an experiment
	dataRootDir string

	// supervisor or not
	isSupervisor bool

	// batch running config
	isGen                bool
	isGenerateForExeFile bool
)

func main() {
	// Start a node.
	pflag.IntVarP(&shardNum, "shardNum", "S", params.ShardNum, "shardNum is an Integer, which indicates that how many shards are deployed. ")
	pflag.IntVarP(&nodeNum, "nodeNum", "N", params.NodesInShard, "nodeNum is an Integer, which indicates how many nodes of each shard are deployed. ")
	pflag.IntVarP(&shardID, "shardID", "s", 0, "shardID is an Integer, which indicates the ID of the shard to which this node belongs. Value range: [0, shardNum). ")
	pflag.IntVarP(&nodeID, "nodeID", "n", 0, "nodeID is an Integer, which indicates the ID of this node. Value range: [0, nodeNum).")
	pflag.IntVarP(&modID, "modID", "m", 3, "modID is an Integer, which indicates the choice ID of methods / consensuses. Value range: [0, 4), representing [CLPA_Broker, CLPA, Broker, Relay]")
	pflag.StringVarP(&dataRootDir, "dataRootDir", "d", params.ExpDataRootDir, "dataRootDir is a string, which defines the RootDir of the experimental data, including ./log, ./record and ./result")
	pflag.BoolVarP(&isSupervisor, "supervisor", "c", false, "isSupervisor is a bool value, which indicates whether this node is a supervisor.")
	pflag.BoolVarP(&isGen, "gen", "g", false, "isGen is a bool value, which indicates whether to generate a batch file")
	pflag.BoolVarP(&isGenerateForExeFile, "shellForExe", "f", false, "isGenerateForExeFile is a bool value, which is effective only if 'isGen' is true; True to generate for an executable, False for 'go run'. ")
	pflag.Parse()

	if isGen {
		if isGenerateForExeFile {
			// Determine the current operating system.
			// Generate the corresponding .bat file or .sh file based on the detected operating system.
			os := runtime.GOOS
			switch os {
			case "windows":
				build.Exebat_Windows_GenerateBatFile(nodeNum, shardNum, modID, dataRootDir)
			case "darwin":
				build.Exebat_MacOS_GenerateShellFile(nodeNum, shardNum, modID, dataRootDir)
			case "linux":
				build.Exebat_Linux_GenerateShellFile(nodeNum, shardNum, modID, dataRootDir)
			}
		} else {
			// Without determining the operating system.
			// Generate a .bat file or .sh file for running `go run`.
			build.GenerateBatFile(nodeNum, shardNum, modID, dataRootDir)
			build.GenerateShellFile(nodeNum, shardNum, modID, dataRootDir)
		}

		return
	}

	// set the experimental data root dir
	params.DataWrite_path = dataRootDir + "/result/"
	params.LogWrite_path = dataRootDir + "/log/"
	params.DatabaseWrite_path = dataRootDir + "/database/"

	if isSupervisor {
		build.BuildSupervisor(uint64(nodeNum), uint64(shardNum), uint64(modID))
	} else {
		build.BuildNewPbftNode(uint64(nodeID), uint64(nodeNum), uint64(shardID), uint64(shardNum), uint64(modID))
	}
}
