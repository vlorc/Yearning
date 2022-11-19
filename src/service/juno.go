package service

import (
	"Yearning-go/src/lib"
	"Yearning-go/src/proto"
	"Yearning-go/src/proxy"
	"context"
	proxy2 "golang.org/x/net/proxy"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type JunoClientProxy struct {
	client     proto.JunoClient
	dialer     proxy2.Dialer
	proxyAlias string
}

var _ proto.JunoClient = &JunoClientProxy{}

func init() {
	lib.NewJunoClient = NewJunoClient
}

func NewJunoClient(conn *grpc.ClientConn, proxyAlias string) proto.JunoClient {
	client := proto.NewJunoClient(conn)

	conf := ProxyService{}.InfoByAlias(proxyAlias)
	if nil == conf {
		return client
	}

	cli := &JunoClientProxy{
		client:     client,
		dialer: proxy.New(conf.Driver, conf.Url, conf.Username, conf.Password, conf.Secret),
		proxyAlias: proxyAlias,
	}


	return cli
}

func (j *JunoClientProxy) begin(ctx context.Context, in *proto.LibraAuditOrder) *proxy.ProxyServer {
	s := &proxy.ProxyServer{
		Name:   j.proxyAlias,
		Target: net.JoinHostPort(in.Source.Addr, strconv.Itoa(int(in.Source.Port))),
		Host:   "127.0.0.1",
		Dial:   j.dialer,
	}

	s.Run()

	source := *in.Source
	source.Addr = "127.0.0.1"
	port, _ := strconv.Atoi(s.Port)
	source.Port = int32(port)
	in.Source = &source

	return s
}

func (j *JunoClientProxy) OrderDeal(ctx context.Context, in *proto.LibraAuditOrder, opts ...grpc.CallOption) (*proto.RecordSet, error) {
	order := *in
	s := j.begin(ctx, &order)
	defer s.Close()

	return j.client.OrderDeal(ctx, &order, opts...)
}

func (j *JunoClientProxy) OrderDMLExec(ctx context.Context, in *proto.LibraAuditOrder, opts ...grpc.CallOption) (*proto.ExecOrder, error) {
	order := *in
	s := j.begin(ctx, &order)
	defer s.Close()

	return j.client.OrderDMLExec(ctx, &order, opts...)
}

func (j *JunoClientProxy) OrderDDLExec(ctx context.Context, in *proto.LibraAuditOrder, opts ...grpc.CallOption) (*proto.ExecOrder, error) {
	order := *in
	s := j.begin(ctx, &order)
	defer s.Close()

	return j.client.OrderDDLExec(ctx, &order, opts...)
}

func (j *JunoClientProxy) AutoTask(ctx context.Context, in *proto.LibraAuditOrder, opts ...grpc.CallOption) (*proto.Isok, error) {
	order := *in
	s := j.begin(ctx, &order)
	defer s.Close()

	return j.client.AutoTask(ctx, &order, opts...)
}

func (j *JunoClientProxy) Query(ctx context.Context, in *proto.LibraAuditOrder, opts ...grpc.CallOption) (*proto.InsulateWordList, error) {
	order := *in
	s := j.begin(ctx, &order)
	defer s.Close()

	return j.client.Query(ctx, &order, opts...)
}

func (j *JunoClientProxy) KillOsc(ctx context.Context, in *proto.LibraAuditOrder, opts ...grpc.CallOption) (*proto.Isok, error) {
	order := *in
	s := j.begin(ctx, &order)
	defer s.Close()

	return j.client.KillOsc(ctx, &order, opts...)
}

func (j *JunoClientProxy) OverrideConfig(ctx context.Context, in *proto.LibraAuditOrder, opts ...grpc.CallOption) (*proto.Isok, error) {
	order := *in
	s := j.begin(ctx, &order)
	defer s.Close()

	return j.client.OverrideConfig(ctx, &order, opts...)
}
