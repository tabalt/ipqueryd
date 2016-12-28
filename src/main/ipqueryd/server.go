package main

import (
	"strings"
	"sync"

	"pb"
)

// import for http
import (
	"encoding/json"
	"html"
	"net/http"

	"github.com/tabalt/gracehttp"
)

// import for grpc
import (
	"net"

	"github.com/tabalt/ipquery"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ipQueryServer struct {
}

func newIpQueryServer() *ipQueryServer {
	return &ipQueryServer{}
}

func (iqs *ipQueryServer) Find(ctx context.Context, params *pb.IpFindRequest) (*pb.IpFindReply, error) {
	result, err := ipquery.Find(params.Ip)
	if err != nil {
		return nil, err
	}

	return &pb.IpFindReply{
		Data: strings.Split(string(result), "\t"),
	}, nil
}

func startServer(conf *Conf) error {
	err := ipquery.Load(conf.DataFile)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	if conf.HttpServerPort != "" {
		wg.Add(1)
		go func() {
			err = startHttpServer(conf.HttpServerPort)
			wg.Done()
		}()
	}

	if conf.GrpcServerPort != "" {
		wg.Add(1)
		go func() {
			err = startGrpcServer(conf.GrpcServerPort)
			wg.Done()
		}()
	}

	wg.Wait()

	return err
}

func startHttpServer(addr string) error {
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/", http.NotFound)

	iqs := newIpQueryServer()
	http.HandleFunc("/find", func(w http.ResponseWriter, r *http.Request) {
		ip := getQueryFromRequest(r, "ip", "")

		result, err := iqs.Find(context.TODO(), &pb.IpFindRequest{Ip: ip})
		if err != nil {
			result = &pb.IpFindReply{Data: nil}
		}
		resp, _ := json.Marshal(result)

		var output []byte
		if cbk := getQueryFromRequest(r, "_callback", ""); cbk != "" {
			lb, rb := []byte(cbk+"("), []byte(");")

			output = make([]byte, 0, (len(lb) + len(resp) + len(rb)))
			output = append(append(append(output, lb...), resp...), rb...)
		} else {
			output = resp
		}

		w.Write(output)
	})

	return gracehttp.ListenAndServe(addr, nil)
}

func startGrpcServer(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterIpQueryServer(grpcServer, newIpQueryServer())

	return grpcServer.Serve(listener)
}

func getQueryFromRequest(r *http.Request, key, defaultValue string) string {
	value := html.EscapeString(r.FormValue(key))
	if value == "" {
		value = defaultValue
	}
	return value
}
