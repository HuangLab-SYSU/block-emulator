This is the second section of the second chapter in the **BlockEmulator** English documentation. 
Firtly, users should download the dataset ([Ethereum-On-chain data](https://xblock.pro/#/dataset/14), published by [xblock](https://xblock.pro/#/)) for **BlockEmulator**.


Parameters refer to variables or values that are used in a system, model, or function to affect its behavior or output. They can be adjusted or tuned to optimize performance or achieve a desired outcome. Parameters can be defined and set at different levels of a system, such as at the software level, hardware level, or algorithm level. In software development, parameters can be used to control program behavior, such as setting the size of a data buffer or the number of iterations in a loop. In scientific research, parameters can be used to adjust experimental conditions or inputs to a model in order to study the effects on the output. Parameters are an important aspect of system design and optimization as they allow for flexibility and customization in achieving desired outcomes.


| Field | Type | Default Value | Explanation|
| --- | --- | -------| --------|
| Block_Interval | int |  5000    |  interval of New block generation |
| MaxBlockSize_global | int |  2000    |  A block contains an upper limit on the number of transactions     |
| InjectSpeed | int |  2000    | transaction injection speed      |
| TotalDataSize | int |  10000    |   The total number of transactions    |
| DataWrite_path | String |  "./result/"     |  The output location of evulation data     |
| LogWrite_path | String |  "./log"      |   The output location of logs    |
| SupervisorAddr | String |  "127.0.0.1:18800"    |    The IP address of Supervisor   |
| FileInput | String |  "../2000000to2999999_BlockTransaction.csv"    |    Original block transaction data location   |


Meanwhile, more specific configuration are shown as follows:

| Field|Type| Default Value| Explanation|
|-----|-----|----------|---------------|
|Init_Balance|big.Int|100000000000000000000000000000000000000000000|The account initial balance|
|NodesinShard|int|4|    The number of node in a shard  |
|ShardNum| int|4|   the number of shared for Blockchain|
|CommitteeMethod|[]string|{"CLPA_Broker", "CLPA", "Broker", "Relay"}| default committee method|
|MeasureBrokerMod|[]string|{"TPS_Broker", "Latency_Broker", "CrossTxRate_Broker", "TxNumberCount_Broker"}| The methods for testing **BrokerChain** Protocol |
|MeasureBrokMeasureRelayModerMod|[]string|{"TPS_Relay", "Latency_Relay", "CrossTxRate_Relay", "TxNumberCount_Relay"}| The methods for testing **Relay** protocol|



