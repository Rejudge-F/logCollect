# kafka-logMgr
# kafka日志收集系统
- 设计目的：简化日志收集
- 使用方法：配置conf/my.conf文件中的COLLECT节点下的配置，设置项为需要收集的日志，然后设置kafka的地址，然后运行程序即可
- 设计原理：使用tailf包追踪文件的变化，然后使用channel来传输对应的message给kafka，然后对应的机器去链接kafka即可消费对应的日志信息
- 提供的kafka测试地址：47.107.54.187:9092
- 启动本程序前需要保证对应的ip地址的kafka和zookeeper已经启动
