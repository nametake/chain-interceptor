package main

import (
	"context"
	"fmt"

	"github.com/nametake/go-func-intercept/pb"
)

func main() {
	pingServer := &PingAPIServer{}
	wrapStruct := &WrapStruct{
		pingServer: pingServer,
	}

	wrapStruct.Call()
	fmt.Println("-------------------")
	wrapStruct.Call(
		printInterceptor,
	)
	fmt.Println("-------------------")
	wrapStruct.Call(
		printInterceptor,
		printInterceptor2,
	)
}

func printInterceptor(ctx context.Context, arg *pb.PingRequest, rpc RPC) (*pb.PingResponse, error) {
	fmt.Println("BEFORE")
	ret, err := rpc(ctx, arg)
	fmt.Println("AFTER")
	return ret, err
}

func printInterceptor2(ctx context.Context, arg *pb.PingRequest, rpc RPC) (*pb.PingResponse, error) {
	fmt.Println("BEFORE2")
	ret, err := rpc(ctx, arg)
	fmt.Println("AFTER2")
	return ret, err
}

type RPC func(context.Context, *pb.PingRequest) (*pb.PingResponse, error)
type Interceptor func(context.Context, *pb.PingRequest, RPC) (*pb.PingResponse, error)

type WrapStruct struct {
	pingServer pb.PingAPIServer
}

func (w *WrapStruct) Call(interceptors ...Interceptor) {
	ctx := context.Background()
	arg := &pb.PingRequest{
		Msg: "PING",
	}

	n := len(interceptors)
	chained := func(ctx context.Context, arg *pb.PingRequest, rpc RPC) (*pb.PingResponse, error) {
		chainer := func(currentInter Interceptor, currentRPC RPC) RPC {
			return func(currentCtx context.Context, currentReq *pb.PingRequest) (*pb.PingResponse, error) {
				return currentInter(currentCtx, currentReq, currentRPC)
			}
		}

		chainedRPC := rpc
		for i := n - 1; i >= 0; i-- {
			chainedRPC = chainer(interceptors[i], chainedRPC)
		}
		return chainedRPC(ctx, arg)
	}

	ret, err := chained(ctx, arg, w.pingServer.Ping)

	fmt.Println(ret, err)
}

type PingAPIServer struct{}

func (p *PingAPIServer) Ping(ctx context.Context, arg *pb.PingRequest) (*pb.PingResponse, error) {
	fmt.Println("CALLED PING")
	return &pb.PingResponse{
		Msg: fmt.Sprintf("PONG: %s", arg.GetMsg()),
	}, nil
}
