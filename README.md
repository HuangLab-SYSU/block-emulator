# Description
This branch, **Fine-tune-lock**, contains the source code for the paper **Account Migration across Blockchain Shards using Fine-tuned Lock Mechanism** (INFOCOM'24). Detailed information about this branch can be found in the paper.

# Start
To start this *Go* program, review the code logic in `./test/test_shard.go`.

Here are some tips:

## Dataset
We have provided a dataset (i.e., `./20W.csv`) in this branch.

## Network configuration
Before running this code, modify the following variables in `./params/config.go`:
- ClientAddr: *ClientAddr* in this branch is similar to *Supervisor* in the `main` branch.
- NodeTable: *NodeTable* in this branch is similar to *IPmap_nodeTable* in the `main` branch, but this map should be filled in manually.

## Run a node
Running a node in this branch is similar to running one in the `main` branch. 
Here is an example:
```
go run main.go -S 4 -f 1 -s S1 -n N1 -t 20W.csv
```
The explanation of this command can be found in `./test/test_shard.go`.

## Batch running (Not supported)
This branch does not support batch running.
However, you can write some code to enable batch running yourself.

# Additional Information
**Earlier version**: This work is built on an earlier version of **BlockEmulator**, so the code structure does not match the current `main` branch.  
**Omitted some security implementations**: Some features related to security or validation, such as Merkle proof and validation, are not implemented in this code. These features are discussed in other papers and are not central to this work.
