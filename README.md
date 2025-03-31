# Description
This branch, **Fine-tune-lock**, contains the source code for the paper **Account Migration across Blockchain Shards using Fine-tuned Lock Mechanism** (INFOCOM'24). Detailed information about this branch can be found in the paper.

# Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Prerequisites

Make sure you have the following software installed on your machine:
* Go: Download and install from [https://golang.org/dl/](https://golang.org/dl/)

## Running the Project

To run the Block Emulator, follow these steps:

1.  **download project:**
   First, clone this project by using the command:
   ```bash
   git clone https://github.com/HuangLab-SYSU/block-emulator.git
   ```
3.  **Build the Executable:**

   Then, create the executable file by using the following command:
   ```bash
   go build -o blockEmulator_Windows_Precompile.exe main.go
   ```
3. **Run the Executable:**
   After creating the executable, run it with the desired options. For example, to generate a .bat file, use the following command:
   
   ```bash
   .\blockEmulator_Windows_Precompile.exe -g -f -S 2 -N 4
   ```
   Here, -g, -f, -S, and -N are options that can be passed to the executable. Each option serves a specific purpose in configuring the behavior of the emulator. For more details about available parameters and their types, please refer to the main branch documentation.
4. **Parameters**
   - -g, --gen            isGen is a bool value, which indicates whether to generate a batch file
   - -n, --nodeID int     nodeID is an Integer, which indicates the ID of this node. Value range: [0, nodeNum). (default -1)
   - -N, --nodeNum int    nodeNum is an Integer, which indicates how many nodes of each shard are deployed.  (default 4)
   - -s, --shardID int    shardID is an Integer, which indicates the ID of the shard to which this node belongs. Value range: [0, shardNum).  (default -1)    
   - -S, --shardNum int   shardNum is an Integer, which indicates that how many shards are deployed.  (default 4)
   - -f, --shellForExe    isGenerateForExeFile is a bool value, which is effective only if 'isGen' is true; True to generate for an executable, False for 'go run'.
   - -c, --supervisor     isSupervisor is a bool value, which indicates whether this node is a supervisor.

# Additional Information
**Omitted some security implementations**: Some features related to security or validation, such as Merkle proof and validation, are not implemented in this code. These features are discussed in other papers and are not central to this work.
