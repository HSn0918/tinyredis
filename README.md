# tiny-redis
## tiny-redis 是一个由 Go 编写的高性能独立缓存服务器。它实现了完整的 RESP（Redis 序列化协议），因此支持所有的 Redis 客户端。
## 特点
*支持基于 RESP 协议的所有客户端。
*支持字符串、列表、集合、哈希数据类型。
*支持 TTL（键-值对将在 TTL 后被删除）。
*完全内存存储。
*支持一些需要的原子操作命令（如 INCR、DECR、INCRBY、MSET、SMOVE 等）。
## 使用
go1.20+ 
```bash
$ go build -o tiny-redis main.go
```
启动tiny-redis服务:
```bash
$ ./tiny-redis
```
使用启动选项命令或配置文件来更改默认设置：
```bash 
$ ./tiny-redis -h
Usage of ./tiny-redis:
  -config string
        Appoint a config file: such as /etc/redis.conf
  -host string
        Bind host ip: default is 127.0.0.1 (default "127.0.0.1")
  -logdir string
        Set log directory: default is /tmp (default "./")
  -loglevel string
        Set log level: default is info (default "info")
  -port int
        Bind a listening port: default is 6379 (default 6379)
```
## 任何 Redis 客户端都可以与 tiny-redis 服务器通信。
例如，可以使用 redis-cli 与 tiny-redis 服务器通信：
```bash
# start a tiny-redis server listening at 12345 port
$ ./tiny-redis 
[info][server.go:25] 2023/09/17 00:55:35 [Server Listen at 127.0.0.1:6379]
[info][server.go:35] 2023/09/17 00:55:40 [127.0.0.1:7810  connected]
$ redis-cli
127.0.0.1:6379> PING
PONG
127.0.0.1:6379> MSET key1 a key2 b
OK
127.0.0.1:6379> MGET key1 key2 nonekey
1) "a"
2) "b"
3) ""
127.0.0.1:6379> RPUSH list1 1 2 3 4 5
(integer) 5
127.0.0.1:6379> LRANGE list1 0 -1
1) "1"
2) "2"
3) "3"
4) "4"
5) "5"
127.0.0.1:6379> TYPE list1
list
127.0.0.1:6379> EXPIRE list1 100
(integer) 1
127.0.0.1:6379> TTL list1
(integer) 96
127.0.0.1:6379> PERSIST list1
(integer) 1
127.0.0.1:6379> TTL list1
(integer) -1

```
## 性能基准测试
性能基准测试的结果是基于 redis-benchmark 工具进行的。[redis-benchmark](https://redis.io/topics/benchmarks)
测试在Lenovo Legion R70002021, 
AMD Ryzen 5 5600H with Redeon Graphics, 
NVIDA Geforce RTX3050 laptop CPU 4GB, 
16GB (3200MHX),
Windows11 with Ubuntu 20.04.6 LTS(WSL2)
`redis-benchmark -c 50-n 200000 -t get`

```text
get: 146716.22 requests per second
set: 153433.08 requests per second
incr: 144334.86 requests per second
lpush: 145313.64 requests per second
rpush: 139470.00 requests per second
lpop: 152226.30 requests per second
rpop: 147929.08 requests per second
sadd: 160599.60 requests per second
hset: 147765.06 requests per second
spop: 144109.50 requests per second

lrange_100: 83880.90 requests per second
lrange_300: 50652.36 requests per second
lrange_500: 37703.82 requests per second
lrange_600: 27895.92 requests per second

mset: 126196.26 requests per second
```
## 可用命令
All commands used as [redis commands](https://redis.io/commands/). You can use any redis client to communicate with thinRedis.

| key     | string      | list   | set         | hash         |
|---------|-------------|--------|-------------|--------------|
| del     | set         | llen   | sadd        | hdel         |
| exists  | get         | lindex | scard       | hexists      |
| keys    | getrange    | lpos   | sdiff       | hget         |
| expire  | setrange    | lpop   | sdirrstore  | hgetall      |
| persist | mget        | rpop   | sinter      | hincrby      |
| ttl     | mset        | lpush  | sinterstore | hincrbyfloat |
| type    | setex       | lpushx | sismember   | hkeys        |
| rename  | setnx       | rpush  | smembers    | hlen         |
|         | strlen      | rpushx | smove       | hmget        |
|      | incr        | lset   | spop        | hset         |
|      | incrby      | lrem   | srandmember | hsetnx       |
|      | decr        | ltrim  | srem        | hvals        |
|      | decrby      | lrange | sunion      | hstrlen      |
|      | incrbyfloat | lmove  | sunionstore | hrandfield   |
|      | append      |        |             |              |
## 文件目录

### 一级目录

```go
|-- README.md
|-- RESP//RESP协议数据类型和解析RESP函数
|-- config//配置
|-- go.mod
|-- logger//日志
|-- main.go
|-- memdb//数据库数据类型
|-- server//服务端
`-- util
```

### 二级目录

```go
|-- README.md
|-- RESP
|   |-- arraydata.go
|   |-- bulkdata.go
|   |-- errordata.go
|   |-- intdata.go
|   |-- parser_test.go
|   |-- parsestream.go
|   |-- plaindata.go
|   |-- stringdata.go
|   `-- structure.go//解析RESP协议
|-- config
|   `-- config.go
|-- go.mod
|-- logger
|   |-- level.go
|   `-- logger.go
|-- main.go
|-- memdb
|   |-- command.go//方法注册函数
|   |-- concurrentmap.go//ConcurrentMap
|   |-- db.go//内存数据库
|   |-- dblock.go//锁
|   |-- hash.go
|   |-- hash_struct.go
|   |-- keys.go
|   |-- keys_test.go
|   |-- list.go
|   |-- list_struct.go
|   |-- list_test.go
|   |-- set.go
|   |-- set_struct.go
|   |-- string.go
|   |-- string_test.go
|   |-- zset.go
|   `-- zset_struct.go
|-- server
|   |-- handler.go//监听端口
|   `-- server.go//启动服务
`-- util
    `-- util.go//hash函数和正则实现
```

