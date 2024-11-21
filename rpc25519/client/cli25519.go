package main

import (
	"flag"
	"fmt"
	stdlog "log"
	"reflect"
	//"net"
	//"net/rpc"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/glycerine/rpc25519"
	//codec "github.com/mars9/codec"
	"github.com/rpcxio/rpcx-benchmark/proto"
	"github.com/rpcxio/rpcx-benchmark/stat"
	"github.com/smallnest/rpcx/log"
	"go.uber.org/ratelimit"
)

var (
	concurrency = flag.Int("c", 1, "concurrency")
	total       = flag.Int("n", 1, "total requests for all clients")
	host        = flag.String("s", "127.0.0.1:8972", "server ip and port")
	pool        = flag.Int("pool", 10, " shared grpc clients")
	rate        = flag.Int("r", 0, "throughputs")
)

/*
func mainOrig() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile) // Add Lshortfile for short file names
	//log.SetLogger(log.NewDefaultLogger(os.Stdout, "", stdlog.LstdFlags|stdlog.Lshortfile, log.LvInfo))

	var rl ratelimit.Limiter
	if *rate > 0 {
		rl = ratelimit.New(*rate)
	}
	// number of concurrent goroutines. Simulate client.
	n := *concurrency
	// number of requests each client needs to send.
	m := *total / n
	log.Printf("concurrency: %d\nrequests per client: %d\n\n", n, m)

	name := "Hello.Say"

	args := proto.PrepareArgs()

	// parameter size
	b := make([]byte, 1024)
	i, _ := args.MarshalTo(b)
	log.Printf("message size: %d bytes\n\n", i)

	// Wait for all tests to complete.
	var wg sync.WaitGroup
	wg.Add(n * m)

	// Total number of requests.
	var trans uint64
	// Return the total number of responses that were normal.
	var transOK uint64

	// Time consumption for each goroutine.
	d := make([][]int64, n, n)

	// create a client connection pool
	var clientIndex uint64
	poolClients := make([]*rpc.Client, 0, *pool)
	for i := 0; i < *pool; i++ {
		conn, err := net.Dial("tcp", *host)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		c := codec.NewClientCodec(conn)
		client := rpc.NewClientWithCodec(c)

		// warmup
		var reply proto.BenchmarkMessage
		for j := 0; j < 5; j++ {
			client.Call(name, args, &reply)
		}
		poolClients = append(poolClients, client)
	}

	// Fence: control the client to start testing at the same time.
	var startWg sync.WaitGroup
	startWg.Add(n + 1) // +1 is because there is a goroutine to record the start time.

	// create a client goroutine and test it.
	startTime := time.Now().UnixNano()
	go func() {
		startWg.Done()
		startWg.Wait()
		startTime = time.Now().UnixNano()
	}()

	for i := 0; i < n; i++ {
		dt := make([]int64, 0, m)
		d = append(d, dt)

		go func(i int) {
			var reply proto.BenchmarkMessage
			startWg.Done()
			startWg.Wait()

			for j := 0; j < m; j++ {
				// Current limiting, here the current waiting time is not counted into the waiting time.
				if rl != nil {
					rl.Take()
				}

				t := time.Now().UnixNano()
				ci := atomic.AddUint64(&clientIndex, 1)
				ci = ci % uint64(*pool)
				client := poolClients[int(ci)]
				err := client.Call(name, args, &reply)
				t = time.Now().UnixNano() - t

				d[i] = append(d[i], t)

				if err == nil && reply.Field1 == "OK" {
					atomic.AddUint64(&transOK, 1)
				}
				// if err != nil {
				// 	log.Error(err)
				// }

				atomic.AddUint64(&trans, 1)
				wg.Done()
			}
		}(i)

	}

	wg.Wait()

	// statistics
	stat.Stats(startTime, *total, d, trans, transOK)
}
*/

func main() {
	flag.Parse()

	//log.SetFlags(log.LstdFlags | log.Lshortfile) // Add Lshortfile for short file names
	log.SetLogger(log.NewDefaultLogger(os.Stdout, "", stdlog.LstdFlags|stdlog.Lshortfile, log.LvInfo))

	var rl ratelimit.Limiter
	if *rate > 0 {
		rl = ratelimit.New(*rate)
	}
	// number of concurrent goroutines. Simulate client.
	n := *concurrency
	// number of requests each client needs to send.
	m := *total / n
	log.Infof("concurrency: %d\nrequests per client: %d\n\n", n, m)

	name := "Hello.Say"

	// parameter size
	b := make([]byte, 1024)

	argsProto := proto.PrepareArgs()
	//fmt.Printf("argsProto = '%#v'\n", argsProto)
	i, _ := argsProto.MarshalTo(b)
	log.Infof("proto message size: %d bytes\n\n", i)
	_ = argsProto

	args := prepareArgs()
	//fmt.Printf("args = '%#v'\n", args)
	b2 := make([]byte, 4096)
	b2, err := args.MarshalMsg(b2[:0])
	if err != nil {
		panic("could not MarshalMsg")
	}
	log.Infof("rpc25519/greenpack message size: %d bytes\n\n", len(b2))

	// Wait for all tests to complete.
	var wg sync.WaitGroup
	wg.Add(n * m)

	// Total number of requests.
	var trans uint64
	// Return the total number of responses that were normal.
	var transOK uint64

	// Time consumption for each goroutine.
	d := make([][]int64, n, n)

	cfg := rpc25519.NewConfig()
	cfg.ClientDialToHostPort = *host
	cfg.TCPonly_no_TLS = true
	cfg.SkipVerifyKeys = true

	// create a client connection pool
	var clientIndex uint64
	poolClients := make([]*rpc25519.Client, 0, *pool)
	for i := 0; i < *pool; i++ {

		cli, err := rpc25519.NewClient(fmt.Sprintf("cli_%v", i), cfg)
		if err != nil {
			log.Infof("bad client config: '%v'\n", err)
			os.Exit(1)
		}
		err = cli.Start()
		if err != nil {
			log.Infof("client could not connect: '%v'\n", err)
			os.Exit(1)
		}
		defer cli.Close()
		//log.Infof("pool client %v connected from local addr='%v'\n", i, cli.LocalAddr())

		/*		conn, err := net.Dial("tcp", *host)
				if err != nil {
					log.Fatal("dialing:", err)
				}
				c := codec.NewClientCodec(conn)
				client := rpc.NewClientWithCodec(c)
		*/

		// warmup
		var reply rpc25519.BenchmarkMessage
		for j := 0; j < 5; j++ {
			cli.Call(name, args, &reply, nil)
			if j == 0 && i == 0 {
				//fmt.Printf("reply = '%#v'\n", reply)
			}
		}
		poolClients = append(poolClients, cli)
	}

	// Fence: control the client to start testing at the same time.
	var startWg sync.WaitGroup
	startWg.Add(n + 1) // +1 is because there is a goroutine to record the start time.

	// create a client goroutine and test it.
	startTime := time.Now().UnixNano()
	go func() {
		startWg.Done()
		startWg.Wait()
		startTime = time.Now().UnixNano()
	}()

	for i := 0; i < n; i++ {
		dt := make([]int64, 0, m)
		d = append(d, dt)

		go func(i int) {
			var reply rpc25519.BenchmarkMessage
			startWg.Done()
			startWg.Wait()

			for j := 0; j < m; j++ {
				// Current limiting, here the current waiting time is not counted into the waiting time.
				if rl != nil {
					rl.Take()
				}

				t := time.Now().UnixNano()
				ci := atomic.AddUint64(&clientIndex, 1)
				ci = ci % uint64(*pool)
				client := poolClients[int(ci)]
				err := client.Call(name, args, &reply, nil)
				t = time.Now().UnixNano() - t

				d[i] = append(d[i], t)

				if err == nil && reply.Field1 == "OK" {
					atomic.AddUint64(&transOK, 1)
				}
				// if err != nil {
				// 	log.Error(err)
				// }

				atomic.AddUint64(&trans, 1)
				wg.Done()
			}
		}(i)

	}

	wg.Wait()

	// statistics
	stat.Stats(startTime, *total, d, trans, transOK)
}

func prepareArgs() *rpc25519.BenchmarkMessage {
	b := true
	var i int32 = 100000
	var s = "许多往事在眼前一幕一幕，变的那麼模糊"

	var args rpc25519.BenchmarkMessage

	v := reflect.ValueOf(&args).Elem()
	num := v.NumField()
	for k := 0; k < num; k++ {
		field := v.Field(k)
		if field.Type().Kind() == reflect.Ptr {
			switch v.Field(k).Type().Elem().Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				field.Set(reflect.ValueOf(&i))
			case reflect.Bool:
				field.Set(reflect.ValueOf(&b))
			case reflect.String:
				field.Set(reflect.ValueOf(&s))
			}
		} else {
			switch field.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				field.SetInt(100000)
			case reflect.Bool:
				field.SetBool(true)
			case reflect.String:
				field.SetString(s)
			}
		}

	}
	return &args
}
