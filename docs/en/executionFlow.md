- It is the fourth section of the second chapter of the **BlockEmulator** English introduction document.

The operation cycle of the entire BlockEmulator can be divided into three stages: **parameter configuration before operation**, **execution**, and **data collection during operation and after operation**

<div algin=center><img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/workflow.png" width=800)></div>
 <center> The period of exection of BlockEmulator </center>


## 1. parameter configuration before operation

This step refers to the specific configuration of **BlockEmulator** in the first section of Chapter 2, where there is a detailed configuration tutorial.

## 2. Execution

- Start each node and the master node reads the shared transaction from the file.
  
- Every period of time, the master node injects a certain number of transactions into the transaction pool.
  
- Every period of time, the master node extracts a certain number of transactions from the transaction pool, packages the block, initiates the PBFT consensus in the shard, and uploads the block to the chain after the consensus ends, updates the state tree, and stores it persistently; While the block is on the chain, the block, transaction and other information are recorded in a CSV file to facilitate subsequent analysis.
  
- For cross-shard transactions, the master node periodically takes the relay transaction from the relay transaction pool and sends it to the master node of the destination shard, and the destination shard master node puts it into the transaction pool and waits to be packaged together with the intra-slice transaction.
  
- After all transactions are executed and put on the chain, the node stops running (a stop operation message is sent by one client).
  
## 3. data collection
   
For this step, refer to the indicators in Section 6 of Chapter 2 and the logs in Section 7 of Chapter 2. In those two sections, there are detailed descriptions of test indicators and descriptions of log generation.
