package lib

import (
	"Yearning-go/src/model"
	pb "Yearning-go/src/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	globalGRPCconns *grpc.ClientConn
	lock            sync.Mutex
	config          atomic.Value
)

var NewJunoClient = func(conn *grpc.ClientConn, proxyAlias string) pb.JunoClient {
	return pb.NewJunoClient(conn)
}

const TIMEOUT = time.Second * 180

func FetchGRPCConn() (*grpc.ClientConn, error) {
	if c := config.Load(); c != nil {
		if c.(*grpc.ClientConn).GetState() == connectivity.Ready {
			return c.(*grpc.ClientConn), nil
		}
	}

	lock.Lock()

	defer lock.Unlock()

	cli, err := newGrpcConn()

	cli.Target()

	if err != nil {
		return nil, err
	}

	config.Store(cli)

	return cli, nil
}

func newGrpcConn() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		model.Grpc,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func TsClient(order *pb.LibraAuditOrder, proxyAlias string) ([]*pb.Record, error) {

	conn, err := FetchGRPCConn()

	if err != nil {
		return nil, err
	}

	c := NewJunoClient(conn, proxyAlias)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer func() {
		cancel()
	}()
	r, err := c.OrderDeal(ctx, order)
	if err != nil {
		log.Printf("could not connect: %v", err)
		return []*pb.Record{}, err
	}
	return r.Record, nil
}

func ExDDLClient(order *pb.LibraAuditOrder, proxyAlias string) {
	// Set up a connection to the server.

	conn, err := FetchGRPCConn()

	if err != nil {
		log.Println(err.Error())
		return
	}

	c := NewJunoClient(conn, proxyAlias)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer func() {
		cancel()
	}()
	_, err = c.OrderDDLExec(ctx, order)
	if err != nil {
		log.Printf("could not connect: %v", err)
		MessagePush(order.WorkId, EVENT_ORDER_EXEC_FAILED, "")
		return
	}
	MessagePush(order.WorkId, EVENT_ORDER_EXEC_SUCCESS, "")
}

func ExDMLClient(order *pb.LibraAuditOrder, proxyAlias string) {

	conn, err := FetchGRPCConn()

	if err != nil {
		log.Println(err.Error())
		return
	}

	// Set up a connection to the server.
	c := NewJunoClient(conn, proxyAlias)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer func() {
		cancel()
	}()
	_, err = c.OrderDMLExec(ctx, order)
	if err != nil {
		log.Printf("could not connect: %v", err)
		MessagePush(order.WorkId, EVENT_ORDER_EXEC_FAILED, "")
		return
	}
	MessagePush(order.WorkId, EVENT_ORDER_EXEC_SUCCESS, "")
}

func ExAutoTask(order *pb.LibraAuditOrder, proxyAlias string) bool {

	conn, err := FetchGRPCConn()

	if err != nil {
		log.Println(err.Error())
		return false
	}

	c := NewJunoClient(conn, proxyAlias)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)

	defer func() {
		cancel()
	}()
	r, err := c.AutoTask(ctx, order)
	if err != nil {
		log.Printf("could not connect: %v", err)
	}
	return r.Ok
}

func ExKillOsc(order *pb.LibraAuditOrder) *pb.Isok {
	conn, err := FetchGRPCConn()

	if err != nil {
		log.Println(err.Error())
	}
	c := NewJunoClient(conn, "")
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer func() {
		cancel()
	}()
	r, err := c.KillOsc(ctx, order)
	if err != nil {
		log.Printf("could not connect: %v", err)
	}
	return r
}

func OverrideConfig(order *pb.LibraAuditOrder) *pb.Isok {
	conn, err := FetchGRPCConn()

	if err != nil {
		log.Println(err.Error())
	}
	c := NewJunoClient(conn, "")
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer func() {
		cancel()
	}()
	r, err := c.OverrideConfig(ctx, order)
	if err != nil {
		log.Printf("could not connect: %v", err)
	}
	return r
}
