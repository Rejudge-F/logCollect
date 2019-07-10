# kafka-logMgr
# kafka日志收集系统
- 设计目的：简化日志收集
- 使用方法：运行tools中的代码，设置etcd中的key，然后运行main的程序，开始收集collect日志
- 设计原理：etcd配置收集信息，使用tailf包追踪文件的变化，然后使用channel来传输对应的message给kafka，然后对应的机器去链接kafka即可消费对应的日志信息
- 提供的kafka测试地址：47.107.54.187:9092
- 启动本程序前需要保证对应的ip地址的kafka和zookeeper以及etcd已经启动
