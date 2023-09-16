# tiny-redis
## tiny-redis 是一个由 Go 编写的高性能独立缓存服务器。它实现了完整的 RESP（Redis 序列化协议），因此支持所有的 Redis 客户端。
## 特点
*支持基于 RESP 协议的所有客户端。
*支持字符串、列表、集合、哈希数据类型。
*支持 TTL（键-值对将在 TTL 后被删除）。
*完全内存存储。
*支持一些需要的原子操作命令（如 INCR、DECR、INCRBY、MSET、SMOVE 等）。
## 使用
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
Usage of ./thinredis:
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
# start a thinRedis server listening at 12345 port
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
Windows11
`redis-benchmark -c 50 -n 200000 -t get`
```text
get: 24452.87 requests per second
set: 25572.18 requests per second
incr: 24055.81 requests per second
lpush: 24218.94 requests per second
rpush: 23245.00 requests per second
lpop: 25371.05 requests per second
rpop: 24654.83 requests per second
sadd: 26766.60 requests per second
hset: 24627.51 requests per second
spop: 24018.25 requests per second

lrange_100: 13980.15 requests per second
lrange_300: 8432.06 requests per second
lrange_500: 6283.97 requests per second
lrange_600: 4649.32 requests per second
mset: 21032.71 requests per second
```
## Support Commands
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

```text

Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t get
====== GET ======
  200000 requests completed in 8.18 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

45.41% <= 1 milliseconds
99.90% <= 2 milliseconds
100.00% <= 3 milliseconds
100.00% <= 3 milliseconds
24452.87 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t set
====== SET ======
  200000 requests completed in 7.82 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

55.07% <= 1 milliseconds
99.84% <= 2 milliseconds
99.97% <= 3 milliseconds
99.99% <= 4 milliseconds
100.00% <= 5 milliseconds
100.00% <= 6 milliseconds
25572.18 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t incr
====== INCR ======
  200000 requests completed in 8.31 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

42.87% <= 1 milliseconds
99.85% <= 2 milliseconds
100.00% <= 3 milliseconds
100.00% <= 4 milliseconds
24055.81 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t lpush
====== LPUSH ======
  200000 requests completed in 8.26 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

46.07% <= 1 milliseconds
99.72% <= 2 milliseconds
99.95% <= 3 milliseconds
99.97% <= 4 milliseconds
99.98% <= 5 milliseconds
99.99% <= 6 milliseconds
100.00% <= 7 milliseconds
100.00% <= 8 milliseconds
100.00% <= 8 milliseconds
24218.94 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t rpush
====== RPUSH ======
  200000 requests completed in 8.60 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

39.28% <= 1 milliseconds
99.40% <= 2 milliseconds
99.89% <= 3 milliseconds
99.97% <= 4 milliseconds
100.00% <= 5 milliseconds
100.00% <= 6 milliseconds
23245.00 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t lpop
====== LPOP ======
  200000 requests completed in 7.88 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

54.34% <= 1 milliseconds
99.88% <= 2 milliseconds
100.00% <= 3 milliseconds
25371.05 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t rpop
====== RPOP ======
  200000 requests completed in 8.11 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

50.07% <= 1 milliseconds
99.60% <= 2 milliseconds
99.90% <= 3 milliseconds
99.96% <= 4 milliseconds
99.97% <= 5 milliseconds
99.97% <= 6 milliseconds
99.98% <= 7 milliseconds
99.99% <= 8 milliseconds
99.99% <= 9 milliseconds
100.00% <= 10 milliseconds
100.00% <= 11 milliseconds
24654.83 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t sadd
====== SADD ======
  200000 requests completed in 7.47 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

64.47% <= 1 milliseconds
99.85% <= 2 milliseconds
99.99% <= 3 milliseconds
99.99% <= 4 milliseconds
100.00% <= 5 milliseconds
26766.60 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t hset
====== HSET ======
  200000 requests completed in 8.12 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

48.91% <= 1 milliseconds
99.79% <= 2 milliseconds
100.00% <= 3 milliseconds
100.00% <= 3 milliseconds
24627.51 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t spop
====== SPOP ======
  200000 requests completed in 8.33 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

44.41% <= 1 milliseconds
99.72% <= 2 milliseconds
99.95% <= 3 milliseconds
99.96% <= 4 milliseconds
99.97% <= 5 milliseconds
99.99% <= 6 milliseconds
100.00% <= 7 milliseconds
100.00% <= 8 milliseconds
100.00% <= 9 milliseconds
24018.25 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t lrange_100
====== LPUSH (needed to benchmark LRANGE) ======
  200000 requests completed in 8.56 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

45.03% <= 1 milliseconds
99.27% <= 2 milliseconds
99.76% <= 3 milliseconds
99.82% <= 4 milliseconds
99.85% <= 5 milliseconds
99.87% <= 6 milliseconds
99.89% <= 7 milliseconds
99.90% <= 9 milliseconds
99.90% <= 10 milliseconds
99.91% <= 11 milliseconds
99.91% <= 12 milliseconds
99.91% <= 13 milliseconds
99.93% <= 14 milliseconds
99.94% <= 15 milliseconds
99.95% <= 16 milliseconds
99.95% <= 17 milliseconds
99.96% <= 21 milliseconds
99.96% <= 25 milliseconds
99.96% <= 26 milliseconds
99.97% <= 27 milliseconds
99.97% <= 28 milliseconds
99.98% <= 76 milliseconds
99.98% <= 77 milliseconds
99.99% <= 158 milliseconds
99.99% <= 159 milliseconds
99.99% <= 161 milliseconds
100.00% <= 162 milliseconds
100.00% <= 163 milliseconds
100.00% <= 163 milliseconds
23350.85 requests per second

====== LRANGE_100 (first 100 elements) ======
  200000 requests completed in 14.31 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

0.63% <= 1 milliseconds
69.08% <= 2 milliseconds
97.73% <= 3 milliseconds
98.38% <= 4 milliseconds
98.80% <= 5 milliseconds
99.17% <= 6 milliseconds
99.44% <= 7 milliseconds
99.59% <= 8 milliseconds
99.70% <= 9 milliseconds
99.78% <= 10 milliseconds
99.83% <= 11 milliseconds
99.85% <= 12 milliseconds
99.87% <= 13 milliseconds
99.89% <= 14 milliseconds
99.91% <= 15 milliseconds
99.93% <= 16 milliseconds
99.95% <= 17 milliseconds
99.97% <= 18 milliseconds
99.98% <= 19 milliseconds
100.00% <= 20 milliseconds
100.00% <= 21 milliseconds
13980.15 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t lrange_300
====== LPUSH (needed to benchmark LRANGE) ======
  200000 requests completed in 7.61 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

62.60% <= 1 milliseconds
99.85% <= 2 milliseconds
99.96% <= 3 milliseconds
99.97% <= 4 milliseconds
99.98% <= 5 milliseconds
99.99% <= 6 milliseconds
100.00% <= 7 milliseconds
100.00% <= 9 milliseconds
100.00% <= 9 milliseconds
26270.85 requests per second

====== LRANGE_300 (first 300 elements) ======
  200000 requests completed in 23.72 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

0.10% <= 1 milliseconds
11.35% <= 2 milliseconds
57.85% <= 3 milliseconds
88.44% <= 4 milliseconds
93.53% <= 5 milliseconds
94.58% <= 6 milliseconds
95.55% <= 7 milliseconds
96.50% <= 8 milliseconds
97.44% <= 9 milliseconds
98.23% <= 10 milliseconds
98.75% <= 11 milliseconds
99.10% <= 12 milliseconds
99.29% <= 13 milliseconds
99.41% <= 14 milliseconds
99.47% <= 15 milliseconds
99.52% <= 16 milliseconds
99.55% <= 17 milliseconds
99.57% <= 18 milliseconds
99.61% <= 19 milliseconds
99.63% <= 20 milliseconds
99.65% <= 21 milliseconds
99.67% <= 22 milliseconds
99.69% <= 23 milliseconds
99.70% <= 24 milliseconds
99.72% <= 25 milliseconds
99.75% <= 26 milliseconds
99.78% <= 27 milliseconds
99.80% <= 28 milliseconds
99.83% <= 29 milliseconds
99.84% <= 30 milliseconds
99.86% <= 31 milliseconds
99.87% <= 32 milliseconds
99.89% <= 33 milliseconds
99.92% <= 34 milliseconds
99.94% <= 35 milliseconds
99.96% <= 36 milliseconds
99.98% <= 37 milliseconds
99.99% <= 38 milliseconds
100.00% <= 39 milliseconds
100.00% <= 40 milliseconds
100.00% <= 41 milliseconds
100.00% <= 44 milliseconds
100.00% <= 46 milliseconds
8432.06 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t lrange_500
====== LPUSH (needed to benchmark LRANGE) ======
  200000 requests completed in 8.25 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

46.63% <= 1 milliseconds
99.67% <= 2 milliseconds
99.95% <= 3 milliseconds
99.97% <= 4 milliseconds
99.98% <= 5 milliseconds
99.98% <= 6 milliseconds
99.98% <= 7 milliseconds
99.99% <= 8 milliseconds
100.00% <= 9 milliseconds
100.00% <= 11 milliseconds
24242.42 requests per second

====== LRANGE_500 (first 450 elements) ======
  200000 requests completed in 31.83 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

0.02% <= 1 milliseconds
1.08% <= 2 milliseconds
22.28% <= 3 milliseconds
57.61% <= 4 milliseconds
81.76% <= 5 milliseconds
90.56% <= 6 milliseconds
92.85% <= 7 milliseconds
93.73% <= 8 milliseconds
94.61% <= 9 milliseconds
95.51% <= 10 milliseconds
96.40% <= 11 milliseconds
97.21% <= 12 milliseconds
97.89% <= 13 milliseconds
98.36% <= 14 milliseconds
98.72% <= 15 milliseconds
98.95% <= 16 milliseconds
99.10% <= 17 milliseconds
99.20% <= 18 milliseconds
99.26% <= 19 milliseconds
99.32% <= 20 milliseconds
99.37% <= 21 milliseconds
99.40% <= 22 milliseconds
99.43% <= 23 milliseconds
99.46% <= 24 milliseconds
99.49% <= 25 milliseconds
99.51% <= 26 milliseconds
99.54% <= 27 milliseconds
99.56% <= 28 milliseconds
99.58% <= 29 milliseconds
99.61% <= 30 milliseconds
99.63% <= 31 milliseconds
99.65% <= 32 milliseconds
99.67% <= 33 milliseconds
99.69% <= 34 milliseconds
99.71% <= 35 milliseconds
99.73% <= 36 milliseconds
99.74% <= 37 milliseconds
99.76% <= 38 milliseconds
99.78% <= 39 milliseconds
99.79% <= 40 milliseconds
99.80% <= 41 milliseconds
99.81% <= 42 milliseconds
99.81% <= 43 milliseconds
99.82% <= 44 milliseconds
99.82% <= 45 milliseconds
99.83% <= 46 milliseconds
99.85% <= 47 milliseconds
99.86% <= 48 milliseconds
99.88% <= 49 milliseconds
99.89% <= 50 milliseconds
99.91% <= 51 milliseconds
99.92% <= 52 milliseconds
99.94% <= 53 milliseconds
99.95% <= 54 milliseconds
99.96% <= 55 milliseconds
99.97% <= 56 milliseconds
99.98% <= 57 milliseconds
99.99% <= 58 milliseconds
99.99% <= 59 milliseconds
100.00% <= 60 milliseconds
100.00% <= 61 milliseconds
100.00% <= 62 milliseconds
100.00% <= 62 milliseconds
6283.97 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t lrange_600
====== LPUSH (needed to benchmark LRANGE) ======
  200000 requests completed in 8.75 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

38.98% <= 1 milliseconds
98.89% <= 2 milliseconds
99.86% <= 3 milliseconds
99.92% <= 4 milliseconds
99.96% <= 5 milliseconds
99.98% <= 6 milliseconds
99.99% <= 7 milliseconds
100.00% <= 8 milliseconds
100.00% <= 8 milliseconds
22846.70 requests per second

====== LRANGE_600 (first 600 elements) ======
  200000 requests completed in 43.02 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

0.03% <= 1 milliseconds
0.28% <= 2 milliseconds
9.18% <= 3 milliseconds
31.02% <= 4 milliseconds
51.34% <= 5 milliseconds
66.89% <= 6 milliseconds
75.31% <= 7 milliseconds
80.19% <= 8 milliseconds
82.11% <= 9 milliseconds
83.28% <= 10 milliseconds
84.65% <= 11 milliseconds
86.23% <= 12 milliseconds
88.04% <= 13 milliseconds
89.84% <= 14 milliseconds
91.67% <= 15 milliseconds
93.38% <= 16 milliseconds
94.79% <= 17 milliseconds
95.98% <= 18 milliseconds
96.91% <= 19 milliseconds
97.60% <= 20 milliseconds
98.05% <= 21 milliseconds
98.37% <= 22 milliseconds
98.58% <= 23 milliseconds
98.73% <= 24 milliseconds
98.83% <= 25 milliseconds
98.91% <= 26 milliseconds
98.97% <= 27 milliseconds
99.02% <= 28 milliseconds
99.07% <= 29 milliseconds
99.11% <= 30 milliseconds
99.13% <= 31 milliseconds
99.16% <= 32 milliseconds
99.19% <= 33 milliseconds
99.22% <= 34 milliseconds
99.25% <= 35 milliseconds
99.28% <= 36 milliseconds
99.30% <= 37 milliseconds
99.32% <= 38 milliseconds
99.35% <= 39 milliseconds
99.37% <= 40 milliseconds
99.40% <= 41 milliseconds
99.42% <= 42 milliseconds
99.44% <= 43 milliseconds
99.45% <= 44 milliseconds
99.47% <= 45 milliseconds
99.48% <= 46 milliseconds
99.50% <= 47 milliseconds
99.51% <= 48 milliseconds
99.52% <= 49 milliseconds
99.53% <= 50 milliseconds
99.54% <= 51 milliseconds
99.56% <= 52 milliseconds
99.57% <= 53 milliseconds
99.58% <= 54 milliseconds
99.59% <= 55 milliseconds
99.60% <= 56 milliseconds
99.61% <= 57 milliseconds
99.61% <= 58 milliseconds
99.62% <= 59 milliseconds
99.64% <= 60 milliseconds
99.65% <= 61 milliseconds
99.66% <= 62 milliseconds
99.68% <= 63 milliseconds
99.70% <= 64 milliseconds
99.71% <= 65 milliseconds
99.74% <= 66 milliseconds
99.75% <= 67 milliseconds
99.77% <= 68 milliseconds
99.78% <= 69 milliseconds
99.80% <= 70 milliseconds
99.82% <= 71 milliseconds
99.85% <= 72 milliseconds
99.87% <= 73 milliseconds
99.90% <= 74 milliseconds
99.92% <= 75 milliseconds
99.94% <= 76 milliseconds
99.96% <= 77 milliseconds
99.97% <= 78 milliseconds
99.98% <= 79 milliseconds
99.99% <= 80 milliseconds
99.99% <= 81 milliseconds
99.99% <= 83 milliseconds
99.99% <= 84 milliseconds
100.00% <= 85 milliseconds
100.00% <= 86 milliseconds
100.00% <= 88 milliseconds
100.00% <= 89 milliseconds
100.00% <= 90 milliseconds
100.00% <= 90 milliseconds
4649.32 requests per second



Administrator@DESKTOP-V0RVQ2M MINGW64 ~
$ redis-benchmark -c 50 -n 200000 -t mset
====== MSET (10 keys) ======
  200000 requests completed in 9.51 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

22.30% <= 1 milliseconds
98.62% <= 2 milliseconds
99.74% <= 3 milliseconds
99.86% <= 4 milliseconds
99.92% <= 5 milliseconds
99.96% <= 6 milliseconds
99.98% <= 7 milliseconds
100.00% <= 8 milliseconds
100.00% <= 9 milliseconds
21032.71 requests per second
```
