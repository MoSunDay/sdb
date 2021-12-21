## [SDB](https://github.com/yemingfeng/sdb) ：纯 golang 开发、数据结构丰富、持久化的 NoSQL 数据库
------

### 为什么需要 SDB？

试想以下业务场景：

- 计数服务：对内容的点赞、播放等数据进行统计
- 评论服务：发布评论后，查看某个内容的评论列表
- 推荐服务：每个用户有一个包含内容和权重推荐列表

以上几个业务场景，都可以通过 MySQL + Redis 的方式实现。 这里的问题是：MySQL 更多的是充当持久化的能力，Redis 充当的是在线服务的读写能力。

那么只使用 Redis 行不行？ 答案是否定的，因为 Redis 无法保证数据不丢失。

那有没有一种存储能够支持高级的数据结构，并能够将数据进行持久化的呢？

答案是：非常少的。有些数据库要么是支持的数据结构不够丰富，要么是接入成本太高，要么是不可控。

为了解决上述问题，SDB 产生了。

------

### SDB 简单介绍

- 纯 golang 开发，核心代码不超过 1k，代码易读
- 数据结构丰富
    - [string](https://github.com/yemingfeng/sdb/blob/master/internal/pb/protobuf-spec/string.proto)
    - [list](https://github.com/yemingfeng/sdb/blob/master/internal/pb/protobuf-spec/list.proto)
    - [set](https://github.com/yemingfeng/sdb/blob/master/internal/pb/protobuf-spec/set.proto)
    - [sorted set](https://github.com/yemingfeng/sdb/blob/master/internal/pb/protobuf-spec/sorted_set.proto)
    - [bloom filter](https://github.com/yemingfeng/sdb/blob/master/internal/pb/protobuf-spec/bloom_filter.proto)
    - [hyper log log](https://github.com/yemingfeng/sdb/blob/master/internal/pb/protobuf-spec/hyper_log_log.proto)
    - [bitset](https://github.com/yemingfeng/sdb/blob/master/internal/pb/protobuf-spec/bitset.proto)
    - [map](https://github.com/yemingfeng/sdb/blob/master/internal/pb/protobuf-spec/map.proto)
    - [geo hash](https://github.com/yemingfeng/sdb/blob/master/internal/pb/protobuf-spec/geo_hash.proto)
    - [pub sub](https://github.com/yemingfeng/sdb/blob/master/internal/pb/protobuf-spec/pub_sub.proto)
- 持久化
    - 兼容 [pebble](https://github.com/cockroachdb/pebble)
      、[leveldb](https://github.com/syndtr/goleveldb)
      、[badger](https://github.com/dgraph-io/badger) 存储引擎
- 监控
    - 支持 prometheus + grafana 监控方案
- 限流
    - 支持每秒 qps 的限流策略
- 慢查询查看
    - 可查看慢查询的请求，进行分析

------

### 快速使用

#### 启动服务端

```shell
sh ./scripts/quick_start.sh
```

**默认使用 pebble 存储引擎**。启动后，端口会监听 9000 端口

#### 客户端使用

```go
package main

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// 连接服务器
	c := pb.NewSDBClient(conn)
	setResponse, err := c.Set(context.Background(),
		&pb.SetRequest{Key: []byte("hello"), Value: []byte("world")})
	log.Printf("setResponse: %+v, err: %+v", setResponse, err)
	getResponse, err := c.Get(context.Background(),
		&pb.GetRequest{Key: []byte("hello")})
	log.Printf("getResponse: %+v, err: %+v", getResponse, err)
}
```

------

### 性能测试

测试脚本：[benchmark](https://github.com/yemingfeng/sdb/blob/master/examples/benchmark_sdb.go)

测试机器：MacBook Pro (13-inch, 2016, Four Thunderbolt 3 Ports)

处理器：2.9GHz 双核 Core i5

内存：8GB

**测试结果： peek QPS > 12k，avg QPS > 7k，set avg time < 70ms，get avg time <
0.2ms**

<img alt="benchmark" src="https://github.com/yemingfeng/sdb/raw/master/docs/benchmark.png" width=80% />

------

### [更多思考](https://github.com/yemingfeng/sdb/wiki/%E9%A6%96%E9%A1%B5)

### 感谢

**感谢开源的力量，这里就不一一列举了，请大家移步 [go.mod](https://github.com/yemingfeng/sdb/blob/master/go.mod)**