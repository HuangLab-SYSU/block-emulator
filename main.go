package main

import (
	"blockEmulator/build"
	"runtime"

	"github.com/spf13/pflag"
)

var (
	shardNum             int
	nodeNum              int
	shardID              int
	nodeID               int
	modID                int
	isSuupervisor        bool
	isGen                bool
	isGenerateForExeFile bool
)

func main() {
	pflag.IntVarP(&shardNum, "shardNum", "S", 2, "indicate that how many shards are deployed")
	pflag.IntVarP(&nodeNum, "nodeNum", "N", 4, "indicate how many nodes of each shard are deployed")
	pflag.IntVarP(&shardID, "shardID", "s", 0, "id of the shard to which this node belongs, for example, 0")
	pflag.IntVarP(&nodeID, "nodeID", "n", 0, "id of this node, for example, 0")
	pflag.IntVarP(&modID, "modID", "m", 3, "choice Committee Method,for example, 0, [CLPA_Broker,CLPA,Broker,Relay]")
	pflag.BoolVarP(&isSuupervisor, "supervisor", "c", false, "whether this node is a supervisor")
	pflag.BoolVarP(&isGen, "gen", "g", false, "generation bat")
	pflag.BoolVarP(&isGenerateForExeFile, "shellForExe", "f", false, "judge whether to generate a batch file for a pre-compiled executable file")
	pflag.Parse()

	if isGen {
		if isGenerateForExeFile {
			// Determine the current operating system.
			// Generate the corresponding .bat file or .sh file based on the detected operating system.
			os := runtime.GOOS
			switch os {
			case "windows":
				build.Exebat_Windows_GenerateBatFile(nodeNum, shardNum, modID)
			case "darwin":
				build.Exebat_MacOS_GenerateShellFile(nodeNum, shardNum, modID)
			case "linux":
				build.Exebat_Linux_GenerateShellFile(nodeNum, shardNum, modID)
			}
		} else {
			// Without determining the operating system.
			// Generate a .bat file or .sh file for running `go run`.
			build.GenerateBatFile(nodeNum, shardNum, modID)
			build.GenerateShellFile(nodeNum, shardNum, modID)
		}

		return
	}

	if isSuupervisor {
		build.BuildSupervisor(uint64(nodeNum), uint64(shardNum), uint64(modID))
	} else {
		build.BuildNewPbftNode(uint64(nodeID), uint64(nodeNum), uint64(shardID), uint64(shardNum), uint64(modID))
	}
}
