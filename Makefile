all:
	cd rpc25519/server; go build
	cd rpc25519/client; go build
	cd grpc/server; go build
	cd grpc/client; go build
	cd rpcx/server; go build
	cd rpcx/client; go build
	cd go_stdrpc/server; go build
	cd go_stdrpc/client; go build


