package main

// startup.go: measure start-up latency time rather than
// steady-state throughput.

import (
	"flag"
	"fmt"
	stdlog "log"
	"os"
	"reflect"
	"time"

	"github.com/glycerine/rpc25519"

	"github.com/glycerine/rpcx-benchmark/proto"
	"github.com/glycerine/rpcx-benchmark/stat"
	"github.com/smallnest/rpcx/log"
)

var (
	total   = flag.Int("n", 1, "total requests")
	host    = flag.String("s", "127.0.0.1:8972", "server ip and port")
	pskPath = flag.String("psk", "", "path to pre shared key")
)

func main() {
	flag.Parse()

	//log.SetFlags(log.LstdFlags | log.Lshortfile) // Add Lshortfile for short file names
	log.SetLogger(log.NewDefaultLogger(os.Stdout, "", stdlog.LstdFlags|stdlog.Lshortfile, log.LvInfo))

	// number of calls.
	n := *total

	log.Infof("calling with new client sequentially %v times to measure start-up time latency.", n)

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

	// Total number of requests.
	var trans uint64
	// Return the total number of responses that were normal.
	var transOK uint64

	// Time consumption for call
	d := make([][]int64, n, n)
	m := 1
	for i := 0; i < n; i++ {
		dt := make([]int64, 0, m)
		d = append(d, dt)
	}

	cfg := rpc25519.NewConfig()
	cfg.ClientDialToHostPort = *host
	cfg.TCPonly_no_TLS = false
	cfg.UseQUIC = false
	cfg.SkipVerifyKeys = false // true
	cfg.PreSharedKeyPath = *pskPath

	startTime := time.Now().UnixNano()

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {

			cli, err := rpc25519.NewClient(fmt.Sprintf("cli_%v", i), cfg)
			if err != nil {
				log.Infof("bad client config: '%v'\n", err)
				os.Exit(1)
			}
			t := time.Now().UnixNano()
			err = cli.Start()
			if err != nil {
				log.Infof("client could not connect: '%v'\n", err)
				os.Exit(1)
			}

			//log.Infof("pool client %v connected from local addr='%v'\n", i, cli.LocalAddr())

			// run one call.
			var reply rpc25519.BenchmarkMessage
			cli.Call(name, args, &reply, nil)

			transOK++
			trans++

			t = time.Now().UnixNano() - t
			d[i] = append(d[i], t)
			cli.Close()
		}
	}

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
