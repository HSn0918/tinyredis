# tiny-redis
[中文](./README_CN.md)
[EN](./README.md)
## Introduction

`tiny-redis` is a high-performance, standalone cache server written in Go that supports persistence. It fully implements the RESP (Redis Serialization Protocol), making it compatible with any Redis client.

## Features

- Supports all clients that use the RESP protocol.
- Provides support for data types such as strings, lists, sets, and hashes.
- Implements TTL (keys will expire after a specified time).
- Fully in-memory storage.
- Supports various atomic commands like `INCR`, `DECR`, `INCRBY`, `MSET`, `SMOVE`, etc.

## Quick Start with Docker

### Building the Docker Image

```bash
$ git clone https://github.com/HSn0918/tinyredis
$ cd tinyredis
$ docker build -t tiny-redis:0.1 .
```

### Running the Server

We have included the `redis-cli` command-line tool in the `/data` directory for easy debugging and usage.

Note: Although the project currently does not support data persistence, mounting the `/data` directory will still provide access to logs and `tiny-redis` and `redis-cli` binaries.

```bash
$ docker run -d \
  --name tiny-redis \
  -p 6379:6379 \
  -v tinyredis-data:/data \
  tiny-redis:0.1
```

## Building from Source

Requires Go 1.20+

```bash
$ go build -o tiny-redis
```

Start the `tiny-redis` server:

```bash
$ ./tiny-redis
```

Use command-line options or configuration files to change the default settings:

```bash
$ ./tiny-redis -h
A tiny Redis server

Usage:
  tiny-redis [flags]
  tiny-redis [command]

Available Commands:
  completion  Generate completion script
  help        Help about any command

Flags:
  -c, --config string     Specify a config file: such as /etc/redis.conf
  -h, --help              Help for tiny-redis
  -H, --host string       Bind host IP: default is 127.0.0.1 (default "0.0.0.0")
  -d, --logdir string     Set log directory: default is /tmp (default "./")
  -l, --loglevel string   Set log level: default is info (default "info")
  -p, --port int          Bind a listening port: default is 6379 (default 6379)

Use "tiny-redis [command] --help" for more information about a command.
```

You can also enable auto-completion using Cobra for the current session:

```bash
$ ./tiny-redis completion zsh > _tiny-redis_completion
$ source _tiny-redis_completion
```

## Redis Client Compatibility

Any Redis client can communicate with the `tiny-redis` server.

> Currently, graphical clients like Medis and AnotherRedisDesktopManager are supported. However, note that some server-related information displayed by these clients may not be accurate.

For example, you can use the `redis-cli` to communicate with the `tiny-redis` server:

```bash
# Start a tiny-redis server listening on port 6379
$ ./tiny-redis
[info][server.go:25] 2023/09/17 00:55:35 [Server Listen at 127.0.0.1:6379]
[info][server.go:35] 2023/09/17 00:55:40 [127.0.0.1:7810 connected]
```

Using `redis-cli`:

```bash
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

## Performance Benchmark

Performance benchmarks are based on the `redis-benchmark` tool. You can find more information about `redis-benchmark` [here](https://redis.io/topics/benchmarks).

Test machine:  

- **Model**: Lenovo Legion R70002021  
- **CPU**: AMD Ryzen 5 5600H with Radeon Graphics, NVIDIA GeForce RTX3050  
- **Memory**: 16GB (3200MHZ)  
- **Environment**: Windows 11 with Ubuntu 20.04.6 LTS (WSL2)  

Command: `redis-benchmark -c 50 -n 200000 -t get`

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

## Available Commands

`tiny-redis` supports several Redis-like commands. You can view the complete list of Redis commands [here](https://redis.io/commands/).

| key     | string      | list   | set         | hash         | Set  |
| ------- | ----------- | ------ | ----------- | ------------ | ---- |
| del     | set         | llen   | sadd        | hdel         | zadd |
| exists  | get         | lindex | scard       | hexists      |      |
| keys    | getrange    | lpos   | sdiff       | hget         |      |
| expire  | setrange    | lpop   | sdirrstore  | hgetall      |      |
| persist | mget        | rpop   | sinter      | hincrby      |      |
| ttl     | mset        | lpush  | sinterstore | hincrbyfloat |      |
| type    | setex       | lpushx | sismember   | hkeys        |      |
| rename  | setnx       | rpush  | smembers    | hlen         |      |
| ping    | strlen      | rpushx | smove       | hmget        |      |
| info    | incr        | lset   | spop        | hset         |      |
|         | incrby      | lrem   | srandmember | hsetnx       |      |
|         | decr        | ltrim  | srem        | hvals        |      |
|         | decrby      | lrange | sunion      | hstrlen      |      |
|         | incrbyfloat | lmove  | sunionstore | hrandfield   |      |
|         | append      |        |             |              |      |

## Directory Structure

### First-level directories

```bash
.
|-- Dockerfile
|-- LICENSE
|-- Makefile
|-- README.md
|-- cmd
|-- go.mod
|-- go.sum
|-- main.go
|-- pkg
`-- sh
```

### Second-level directories

```bash
.
|-- Dockerfile
|-- LICENSE
|-- Makefile
|-- README.md
|-- cmd
|   `-- init.go
|-- go.mod
|-- go.sum
|-- main.go
|-- pkg
|   |-- RESP
|   |   |-- arraydata.go
|   |   |-- bulkdata.go
|   |   |-- errordata.go
|   |   |-- intdata.go
|   |   |-- parser_test.go
|   |   |-- parsestream.go
|   |   |-- plaindata.go
|   |   |-- stringdata.go
|   |   `-- structure.go
|   |-- config
|   |   `-- config.go
|   |-- logger
|   |   |-- level.go
|   |   `-- logger.go
|   |-- memdb
|   |   |-- command.go
|   |   |-- concurrentmap.go
|   |   |-- concurrentmap_test.go
|   |   |-- db.go
|   |   |-- dblock.go
|   |   |-- hash.go
|   |   |-- hash_struct.go
|   |   |-- info.go
|   |   |-- keys.go
|   |   |-- keys_test.go
|   |   |-- list.go
|   |   |-- list_struct.go
|   |   |-- list_test.go
|   |   |-- set.go
|   |   |-- set_struct.go
|   |   |-- string.go
|   |   |-- string_test.go
|   |   |-- zset.go
|   |   |-- zset_struct.go
|   |   `-- zset_test.go
|   |-- server
|   |   |-- aof.go
|   |   |-- handler.go
|   |   `-- server.go
|   `-- util
|       `-- util.go
`-- sh
    `-- testAof
```

---

