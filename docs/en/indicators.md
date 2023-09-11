**BlockEmulator** metrics will be measured by the supervisor process. Hence, if you want to call metrics in **BlockEmulator**, you need to specify the name of the supervisor when creating the supervisor:

    ```
    1 spv:= new(supervisor.Supervisor)
    2 spv.NewSupervisor(supervisor_ip, chainConfig, committeeMethod, mearsureModNames...)
    ```

In **Golang**, *mearsureModNames* refers to "variadic parameters"

Currently, **BlockEmulator** supports eight kinds of indicators (four classes):

<table>
    <tr>
        <th></th><th>Type of Indicator</th><th>Explaination</th>
    </tr>
    <tr>
        <td rowspan="4">Mode of Cross-shard Transaction (Relay)</td><td>Throughput Per Seconds (TPS)</td><td>For the throughput test of the Relay method, the throughput is the average number of transactions processed by the blockchain system per second; It is one of the basic metrics that blockchain uses to measure performance</td>
    </tr>
    <tr>
        <td>Transaction Confirmation Latency (TCL)</td><td>For the transaction confirmation latency (TCL) test of the Relay method, the TCL of a transaction refers to the time it takes for the transaction to enter the transaction pool to the final confirmation on the chain, while this indicator tests the average TCL of a bunch of transactions; It is also one of the basic metrics that blockchain uses to evaluate performance</td>
    </tr>
    <tr>
        <td>The proportion of cross-shard transactions (Cross TxRate_relay)</td><td>The proportion of cross-shard transactions for the Relay method refers to the proportion of cross-shard transactions in all transactions processed by the blockchain system; It is often used to evaluate the performance of an account reallocation algorithm</td>
    </tr>
    <tr>
        <td>The total transaction
        (TXNumberCount_Relay)</td><td>Measuring the total number of transactions can be used to count the number of final on-chain transactions.</td>
     <tr>
        <td rowspan="4">Mode of Broker account 
        (Broker account)</td><td>Throughput Per Seconds (TPS)</td><td>For the throughput test of the Relay method, the throughput is the average number of transactions processed by the blockchain system per second; It is one of the basic metrics that blockchain uses to measure performance</td>
    </tr>
    <tr>
   <td>Transaction Confirmation Latency (TCL)</td><td>For the **Broker** account method of transaction confirmation latency (TCL) test, the TCL of a transaction refers to the time it takes for the transaction to enter the transaction pool to the final confirmation on the chain, while this indicator tests the average TCL of a bunch of transactions; It is also one of the basic metrics that blockchain uses to evaluate performance</td>
    </tr>
    <tr>
        <td>The proportion of cross-shard transactions (Cross TxRate_relay)</td><td>The proportion of cross-shard transactions for the **broker** account method refers to the proportion of cross-shard transactions in all transactions processed by the blockchain system; It is often used to evaluate the performance of an account reallocation algorithm</td>
    </tr>
    <tr>
        <td>The total transaction
        (TXNumberCount_Relay)</td><td>Measuring the total number of transactions can be used to count the number of final on-chain transactions.</td>
     <tr>
        <td rowspan="4"></td><td></td><td></td>
    </tr>
</table>

(Note: This version of BlockEmulator will count a cross-shard transaction as two 0.5 transactions. )
The above indicators will be outputed after the execution of **Supervisor** with the follows format:

    ```
    1 var metricByEpoch []float // output the metric result in an array, the index of the metric result is set by epoch
    2 var metricResultAll float // the metric result throughout the entire running
    ```

That is, the measurement results will be placed into the array metricByEpoch one by one according to Epoch, and the measurement results of the entire running process will be written to metricResultAll. 

In addition, the indicators are written to params in the form of .csv files. DataWrite_path/supervisor_measureOutput folder.
