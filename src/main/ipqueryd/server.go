package main

import (
	"net"
	"strings"

	"github.com/tabalt/ipquery"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"pb"
)

import (
	"encoding/json"
	"html"
	"net/http"

	"github.com/tabalt/gracehttp"
)

type ipQueryServer struct {
}

func newIpQueryServer() *ipQueryServer {
	return &ipQueryServer{}
}

func (iqs *ipQueryServer) Find(ctx context.Context, params *pb.FindParams) (*pb.FindResult, error) {
	result, err := ipquery.Find(params.Ip)
	if err != nil {
		return nil, err
	}

	return &pb.FindResult{
		Data: strings.Split(string(result), "\t"),
	}, nil
}

func startServer(conf *Conf) error {
	err := ipquery.Load(conf.DataFile)
	if err != nil {
		return err
	}

	var srvErr error = nil

	if srvErr == nil && conf.HttpServerPort != "" {
		srvErr = startHttpServer(conf.HttpServerPort)
	}

	if srvErr == nil && conf.GrpcServerPort != "" {
		srvErr = startGrpcServer(conf.GrpcServerPort)
	}

	return srvErr
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

func startHttpServer(addr string) error {
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/", http.NotFound)

	iqs := newIpQueryServer()
	http.HandleFunc("/find", func(w http.ResponseWriter, r *http.Request) {
		ip := html.EscapeString(r.FormValue("ip"))

		result, err := iqs.Find(context.TODO(), &pb.FindParams{Ip: ip})
		if err != nil {
			result = &pb.FindResult{Data: nil}
		}

		resp, _ := json.Marshal(result)
		w.Write(resp)
	})

	return gracehttp.ListenAndServe(addr, nil)
}
