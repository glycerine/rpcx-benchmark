# RPC benchmarks including [rpc25519](https://github.com/glycerine/rpc25519).

rpc25519 uses a BenchmarkMessage without pointers,
which is more realistic (when using greenpack serialization) compared 
to protobuf version with pointers everywhere. The messages are 
approximately the same size, with the greenpack version saving 14 bytes: 
the serialized message size is 567 bytes for greenpack/rpc25519; 
versus 581 bytes for protobufs/the others.

Performance wise, `rpc25519` has throughput on par with the Go standard library's 
net/rpc, and has slightly better tail latency. Both rpc25519 and
net/rpc out perform the other rpc packages. rpc25519 has twice
the throughput and half the tail latency of gRPC and rpcx.

This makes some sense as rpc25519 reuses most of the net/rpc code.
`rpc25519` enhances the standard package by changing to greenpack 
instead of gob encoding; providing useful header information;
allowing one-way calls as well as two-way calls; and supporting 
server initiated messages and a simple []byte based API too.
See https://github.com/glycerine/rpc25519 for full details.

Even fully encrypted with TLSv1.3 (see the last benchmark), 
the rpc25519 system is still lighter and faster than any other.

On a 48 core linux box:
====================

rpc25519 (TCP; no encryption like the other rpc systems below)
--------
~~~
~/go/src/github.com/rpcxio/rpcx-benchmark/rpc25519/client (rpc25519) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 07:24:39 cli25519.go:159: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 07:24:39 cli25519.go:179: INFO : rpc25519/greenpack message size: 567 bytes

took 24328 ms for 10000000 requests
sent     requests    : 10000000
received requests    : 10000000
received requests_OK : 10000000
throughput  (TPS)    : 411048

mean: 2411150 ns, median: 1284608 ns, max: 72172134 ns, min: 52199 ns, p99.9: 26249190 ns

mean: 2 ms, median: 1 ms, 

max: 72 ms, min: 0 ms, p99.9: 26 ms
~~~

go_stdrpc (net/rpc from standard lib, with protobuf codec)
---------
~~~
~/go/src/github.com/rpcxio/rpcx-benchmark/go_stdrpc/client (master) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 01:23:56 gostd_client.go:41: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 01:23:56 gostd_client.go:50: INFO : message size: 581 bytes

took 23901 ms for 10000000 requests
sent     requests    : 10000000
received requests    : 10000000
received requests_OK : 10000000
throughput  (TPS)    : 418392
mean: 2361007 ns, median: 725327 ns, max: 77126717 ns, min: 36479 ns, p99.9: 32380807 ns

mean: 2 ms, median: 0 ms, 

max: 77 ms, min: 0 ms, p99.9: 32 ms
~~~


grpc
----
~~~
~/go/src/github.com/rpcxio/rpcx-benchmark/grpc/client (master) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 01:34:26 grpc_client.go:44: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 01:34:26 grpc_client.go:47: INFO : Servers: 127.0.0.1:8972

2024/11/21 01:34:26 grpc_client.go:53: INFO : message size: 581 bytes

took 53981 ms for 10000000 requests
sent     requests    : 10000000
received requests    : 10000000
received requests_OK : 10000000
throughput  (TPS)    : 185250
mean: 5331585 ns, median: 1652778 ns, max: 137626060 ns, min: 87967 ns, p99.9: 67124546 ns

mean: 5 ms, median: 1 ms, 

max: 137 ms, min: 0 ms, p99.9: 67 ms
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

took 46658 ms for 10000000 requests
sent     requests    : 10000000
received requests    : 10000000
received requests_OK : 10000000
throughput  (TPS)    : 214325
mean: 4612061 ns, median: 496472 ns, max: 174467900 ns, min: 38393 ns, p99.9: 61442873 ns

mean: 4 ms, median: 0 ms, 

max: 174 ms, min: 0 ms, p99.9: 61 ms
~~~

rpcx at v1.7.8:
~~~
~/go/src/github.com/rpcxio/rpcx-benchmark/rpcx/client (master) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 01:42:23 rpcx_client.go:43: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 01:42:23 rpcx_client.go:51: INFO : Servers: 127.0.0.1:8972

2024/11/21 01:42:23 rpcx_client.go:62: INFO : message size: 581 bytes

took 37275 ms for 10000000 requests
sent     requests    : 10000000
received requests    : 10000000
received requests_OK : 10000000
throughput  (TPS)    : 268276
mean: 3681227 ns, median: 557679 ns, max: 135436762 ns, min: 43893 ns, p99.9: 54399804 ns

mean: 3 ms, median: 0 ms, 

max: 135 ms, min: 0 ms, p99.9: 54 ms
~~~

when we add full encryption (real world use):
========================================

rpc25519 (over QUIC; fully encrypted with TLS v1.3)
--------
~~~
~/go/src/github.com/glycerine/rpcx-benchmark/rpc25519/client (master) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 08:33:58 cli25519.go:159: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 08:33:58 cli25519.go:179: INFO : rpc25519/greenpack message size: 567 bytes

took 107043 ms for 10000000 requests
sent     requests    : 10000000
received requests    : 10000000
received requests_OK : 10000000
throughput  (TPS)    : 93420
mean: 10656358 ns, median: 7933962 ns, max: 138034211 ns, min: 79631 ns, p99.9: 57847390 ns

mean: 10 ms, median: 7 ms, 

max: 138 ms, min: 0 ms, p99.9: 57 ms
~~~

rpc25519 (over TLS v1.3; so TCP now with full encryption)
--------
~~~
~/go/src/github.com/glycerine/rpcx-benchmark/rpc25519/client (master) $ ./client -n 10000000 -c 1000 -pool 1000
2024/11/21 08:38:26 cli25519.go:159: INFO : concurrency: 1000
requests per client: 10000

2024/11/21 08:38:26 cli25519.go:179: INFO : rpc25519/greenpack message size: 567 bytes

took 25788 ms for 10000000 requests
sent     requests    : 10000000
received requests    : 10000000
received requests_OK : 10000000
throughput  (TPS)    : 387777
mean: 2554761 ns, median: 1323340 ns, max: 65111049 ns, min: 63751 ns, p99.9: 27441302 ns

mean: 2 ms, median: 1 ms, 

max: 65 ms, min: 0 ms, p99.9: 27 ms
~~~

rpc25519 (TCP + ASCON 128a light-weight crypto using a pre-shared key, no TLS)

Hilariously, using ASCON 128a encryption over TCP actually speeds things
up compared to TPC only (which is doing no encryption).
~~~
 ./client -n 10000000 -c 1000 -pool 1000 -psk binary.psk
2024/11/23 12:47:38 cli25519.go:46: INFO : concurrency: 1000
requests per client: 10000

rpc25519/greenpack message size: 567 bytes

took 22631 ms for 10000000 requests
sent     requests    : 10000000
received requests    : 10000000
received requests_OK : 10000000
throughput  (TPS)    : 441871

mean: 2244078 ns, median: 1257987 ns, max: 61322020 ns, min: 59112 ns, p99.9: 24072804 ns

mean: 2 ms, median: 1 ms, 

max: 61 ms, min: 0 ms, p99.9: 24 ms
~~~

rpc25519 with TLSv1.3 + ASCON 128a with pre-shared-key; so
two layers of symmetric cryptography for post-quantum resistance.

~~~
$ ./client -n 10000000 -c 1000 -pool 1000 -psk ../server/certs/node.key 
2024/11/23 13:06:08 cli25519.go:46: INFO : concurrency: 1000
requests per client: 10000

2024/11/23 13:06:08 cli25519.go:56: INFO : proto message size: 581 bytes

2024/11/23 13:06:08 cli25519.go:66: INFO : rpc25519/greenpack message size: 567 bytes

took 23859 ms for 10000000 requests
sent     requests    : 10000000
received requests    : 10000000
received requests_OK : 10000000
throughput  (TPS)    : 419129
mean: 2363927 ns, median: 1313512 ns, max: 65219907 ns, min: 64302 ns, p99.9: 25258699 ns

mean: 2 ms, median: 1 ms, 

max: 65 ms, min: 0 ms, p99.9: 25 ms
~~~

latency of connect and one call roundtrip (same host)
===========================

The above was measuring throughput at steady state.
The measurements, as the rpcx benchmark was constructed,
do not include the connection time or the first five
calls (to "warmup").

If, instead, we are interested in latency: the time
to connect and make a single roundtrip, then the
client/startup.go benchmark is useful.

# Latency for TCP/TLSv1.3:

~~~
$ rpcx-benchmark/rpc25519/client (master) $ ./startup  -n 10000
INFO : calling with new client sequentially 10000 times to measure start-up time latency.
INFO : proto message size: 581 bytes
INFO : rpc25519/greenpack message size: 567 bytes

INFO : took 32076 ms for 10000 requests
INFO : sent     requests    : 10000
INFO : received requests    : 10000
received requests_OK : 10000

INFO : mean: 3191267 ns, median: 3203684 ns, 
       max: 36925381 ns, min: 1395699 ns, p99.9: 10928644 ns
INFO : mean: 3 ms, median: 3 ms, 
       max: 36 ms, min: 1 ms, p99.9: 10 ms
~~~

# Latency for QUIC (UDP):

~~~
$ rpcx-benchmark/rpc25519/client (master) $ ./startup  -n 10000
INFO : calling with new client sequentially 10000 times to measure start-up time latency.
INFO : proto message size: 581 bytes

INFO : rpc25519/greenpack message size: 567 bytes

INFO : took 40481 ms for 10000 requests
INFO : sent     requests    : 10000
INFO : received requests    : 10000
INFO : received requests_OK : 10000

INFO : mean: 3989663 ns, median: 3039277 ns, 
       max: 32667285 ns, min: 2075870 ns, p99.9: 23145255 ns
INFO : mean: 3 ms, median: 3 ms, 
       max: 32 ms, min: 2 ms, p99.9: 23 ms
~~~

