# RPC benchmarks including [rpc25519](https://github.com/glycerine/rpc25519).

On a 48 core linux box:
====================

rpc25519
--------
~~~
~/go/src/github.com/rpcxio/rpcx-benchmark/rpc25519/client (rpc25519) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 07:24:39 cli25519.go:159: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 07:24:39 cli25519.go:179: INFO : rpc25519/greenpack message size: 567 bytes

2024/11/21 07:25:05 stats.go:15: INFO : took 24328 ms for 10000000 requests
2024/11/21 07:25:07 stats.go:36: INFO : sent     requests    : 10000000
2024/11/21 07:25:07 stats.go:37: INFO : received requests    : 10000000
2024/11/21 07:25:07 stats.go:38: INFO : received requests_OK : 10000000
2024/11/21 07:25:07 stats.go:42: INFO : throughput  (TPS)    : 411048

2024/11/21 07:25:07 stats.go:45: INFO : mean: 2411150 ns, median: 1284608 ns, max: 72172134 ns, min: 52199 ns, p99.9: 26249190 ns
2024/11/21 07:25:07 stats.go:46: INFO : mean: 2 ms, median: 1 ms, max: 72 ms, min: 0 ms, p99.9: 26 ms
~~~


go_stdrpc
---------
~~~
~/go/src/github.com/rpcxio/rpcx-benchmark/go_stdrpc/client (rpc25519) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 01:23:56 gostd_client.go:41: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 01:23:56 gostd_client.go:50: INFO : message size: 581 bytes

2024/11/21 01:24:20 stats.go:15: INFO : took 23901 ms for 10000000 requests
2024/11/21 01:24:22 stats.go:36: INFO : sent     requests    : 10000000
2024/11/21 01:24:22 stats.go:37: INFO : received requests    : 10000000
2024/11/21 01:24:22 stats.go:38: INFO : received requests_OK : 10000000
2024/11/21 01:24:22 stats.go:42: INFO : throughput  (TPS)    : 418392

2024/11/21 01:24:22 stats.go:45: INFO : mean: 2361007 ns, median: 725327 ns, max: 77126717 ns, min: 36479 ns, p99.9: 32380807 ns
2024/11/21 01:24:22 stats.go:46: INFO : mean: 2 ms, median: 0 ms, max: 77 ms, min: 0 ms, p99.9: 32 ms
~~~


grpc
----
~~~
~/go/src/github.com/rpcxio/rpcx-benchmark/grpc/client (rpc25519) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 01:34:26 grpc_client.go:44: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 01:34:26 grpc_client.go:47: INFO : Servers: 127.0.0.1:8972

2024/11/21 01:34:26 grpc_client.go:53: INFO : message size: 581 bytes

2024/11/21 01:35:21 stats.go:15: INFO : took 53981 ms for 10000000 requests
2024/11/21 01:35:23 stats.go:36: INFO : sent     requests    : 10000000
2024/11/21 01:35:23 stats.go:37: INFO : received requests    : 10000000
2024/11/21 01:35:23 stats.go:38: INFO : received requests_OK : 10000000
2024/11/21 01:35:23 stats.go:42: INFO : throughput  (TPS)    : 185250

2024/11/21 01:35:23 stats.go:45: INFO : mean: 5331585 ns, median: 1652778 ns, max: 137626060 ns, min: 87967 ns, p99.9: 67124546 ns
2024/11/21 01:35:23 stats.go:46: INFO : mean: 5 ms, median: 1 ms, max: 137 ms, min: 0 ms, p99.9: 67 ms
~~~

rpcx
----
rpcx at v1.8.32:
~~~
~/go/src/github.com/rpcxio/rpcx-benchmark/rpcx/client (master) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 01:51:42 rpcx_client.go:43: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 01:51:42 rpcx_client.go:51: INFO : Servers: 127.0.0.1:8972

2024/11/21 01:51:42 rpcx_client.go:62: INFO : message size: 581 bytes

2024/11/21 01:52:29 stats.go:15: INFO : took 46658 ms for 10000000 requests
2024/11/21 01:52:31 stats.go:36: INFO : sent     requests    : 10000000
2024/11/21 01:52:31 stats.go:37: INFO : received requests    : 10000000
2024/11/21 01:52:31 stats.go:38: INFO : received requests_OK : 10000000
2024/11/21 01:52:31 stats.go:42: INFO : throughput  (TPS)    : 214325

2024/11/21 01:52:31 stats.go:45: INFO : mean: 4612061 ns, median: 496472 ns, max: 174467900 ns, min: 38393 ns, p99.9: 61442873 ns
2024/11/21 01:52:31 stats.go:46: INFO : mean: 4 ms, median: 0 ms, max: 174 ms, min: 0 ms, p99.9: 61 ms
~~~

rpcx at v1.7.8:
~~~
~/go/src/github.com/rpcxio/rpcx-benchmark/rpcx/client (rpc25519) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 01:42:23 rpcx_client.go:43: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 01:42:23 rpcx_client.go:51: INFO : Servers: 127.0.0.1:8972

2024/11/21 01:42:23 rpcx_client.go:62: INFO : message size: 581 bytes

2024/11/21 01:43:01 stats.go:15: INFO : took 37275 ms for 10000000 requests
2024/11/21 01:43:03 stats.go:36: INFO : sent     requests    : 10000000
2024/11/21 01:43:03 stats.go:37: INFO : received requests    : 10000000
2024/11/21 01:43:03 stats.go:38: INFO : received requests_OK : 10000000
2024/11/21 01:43:03 stats.go:42: INFO : throughput  (TPS)    : 268276

2024/11/21 01:43:03 stats.go:45: INFO : mean: 3681227 ns, median: 557679 ns, max: 135436762 ns, min: 43893 ns, p99.9: 54399804 ns
2024/11/21 01:43:03 stats.go:46: INFO : mean: 3 ms, median: 0 ms, max: 135 ms, min: 0 ms, p99.9: 54 ms
~~~
