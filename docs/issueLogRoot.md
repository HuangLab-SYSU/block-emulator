# 问题回复 和 Bug 修复

这个文件作为问题日志的根文件，用来索引一些代码修改的记录和思路。

| 问题标题 (文件链接)                       | 问题描述                                                               | 解决方案                                                     |
| ------------------------------ | ---------------------------------------------------------------------- | ------------------------------------------------------------ |
| [PBFT 共识间隔设置](./IssueLogs/pbft_interval_config_dir/pbft_interval_config.md)              | PBFT 的消息在共识间隔很低的时候，会出现乱序。                          | 使用 Go routine 并行处理消息，实现消息按照预期顺序执行。     |
| [CLPACommitteeModule 出现死循环](./IssueLogs/clpa_dead_loop_dir/clpa_dead_loop.md) | 将分片数目设置为 1 时，Supervisor 节点会在第一次运行 CLPA 算法后卡住。 | 加上判断，分片数目为 1 时，Supervisor 节点不执行 CLPA 算法。 |
| [AddAccount 函数抛出空指针异常](./IssueLogs/addaccount_nil_pointer_dir/addaccount_nil_pointer.md) | 当 AddAccount 函数访问并更新完 MPT 后，如果 MPT 并没有添加新的树节点，那么就会抛出空指针异常。 | 提前检查此次调用函数是否添加了新的树节点。 |