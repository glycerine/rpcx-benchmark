package main

import (
	"flag"
	"fmt"
	//"io"
	"log"
	//"net"
	"net/http"
	_ "net/http/pprof"
	//"net/rpc"
	"runtime"
	"time"

	//codec "github.com/mars9/codec"
	"github.com/glycerine/rpc25519"
	//"github.com/rpcxio/rpcx-benchmark/proto"
)

var _ = runtime.Gosched

//go:generate greenpack

type Hello struct {
	Placeholder int // greenpack will ignore if no public fields
}

func (t *Hello) Say(args *rpc25519.BenchmarkMessage, reply *rpc25519.BenchmarkMessage) error {
	args.Field1 = "OK"
	args.Field2 = 100
	*reply = *args
	if *delay > 0 {
		time.Sleep(*delay)
	} else {
		runtime.Gosched()
	}
	//fmt.Printf("Hello.Say called!\n")
	return nil
}

var (
	host       = flag.String("s", "127.0.0.1:8972", "listened ip and port")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	delay      = flag.Duration("delay", 0, "delay to mock business processing")
	debugAddr  = flag.String("d", "127.0.0.1:9981", "server ip and port")
)

/*
func mainOrigNetRpc() {
	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe(*debugAddr, nil))
	}()

	server := rpc.NewServer()
	server.Register(new(Hello))

	ln, err := net.Listen("tcp", *host)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("rpc.Serve: accept:", err.Error())
			return
		}
		go serveConn(server, conn)
	}
}

func serveConn(server *rpc.Server, conn io.ReadWriteCloser) {
	srv := codec.NewServerCodec(conn)
	server.ServeCodec(srv)
}
*/

func main() {
	flag.Parse()

	go func() {
		log.Printf("pprof web listening at %v\n", *debugAddr)
		log.Println(http.ListenAndServe(*debugAddr, nil))
	}()

	cfg := rpc25519.NewConfig()
	cfg.ServerAddr = *host
	cfg.TCPonly_no_TLS = true
	cfg.UseQUIC = false
	cfg.SkipVerifyKeys = true

	srv := rpc25519.NewServer("srv", cfg)
	defer srv.Close()

	//srv.Register2Func(customEcho)
	srv.Register(new(Hello))

	serverAddr, err := srv.Start()
	if err != nil {
		panic(fmt.Sprintf("could not start rpc25519.Server with config = '%#v'; err='%v'", cfg, err))
	}
	log.Printf("server listening on '%v'\n", serverAddr)
	select {}
}

/*
// echo implements rpc25519.TwoWayFunc
func customEcho(req, reply *rpc25519.Message) error {

	args.Field1 = "OK"
	args.Field2 = 100
	*reply = *args
	if *delay > 0 {
		time.Sleep(*delay)
	} else {
		runtime.Gosched()
	}
	return nil

	//vv("callback to echo: with msg='%#v'", in)
	//reply.JobSerz = append(req.JobSerz, []byte(fmt.Sprintf("\n with time customEcho sees this: '%v'", time.Now()))...)
	//return nil
}
*/