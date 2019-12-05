package main

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/nametake/chain-interceptor/pb"
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

func printInterceptor(ctx context.Context, req proto.Message, rpc func(context.Context, proto.Message) (proto.Message, error)) (proto.Message, error) {
	fmt.Println("BEFORE")
	req, err := rpc(ctx, req)
	fmt.Println("AFTER")
	return req, err
}

func printInterceptor2(ctx context.Context, req proto.Message, rpc func(context.Context, proto.Message) (proto.Message, error)) (proto.Message, error) {
	fmt.Println("BEFORE2")
	req, err := rpc(ctx, req)
	fmt.Println("AFTER2")
	return req, err
}

type WrapStruct struct {
	pingServer pb.PingAPIServer
}

func (w *WrapStruct) Call(interceptors ...func(context.Context, proto.Message, func(context.Context, proto.Message) (proto.Message, error)) (proto.Message, error)) {
	ctx := context.Background()
	req := &pb.PingRequest{
		Msg: "PING",
	}

	n := len(interceptors)

	chained := func(ctx context.Context, req proto.Message, rpc func(context.Context, proto.Message) (proto.Message, error)) (proto.Message, error) {
		chainer := func(
			currentInter func(context.Context, proto.Message, func(context.Context, proto.Message) (proto.Message, error)) (proto.Message, error),
			currentHandler func(context.Context, proto.Message) (proto.Message, error),
		) func(context.Context, proto.Message) (proto.Message, error) {
			return func(currentCtx context.Context, currentReq proto.Message) (proto.Message, error) {
				return currentInter(currentCtx, currentReq, currentHandler)
			}
		}

		chainedRPC := rpc
		for i := n - 1; i >= 0; i-- {
			chainedRPC = chainer(interceptors[i], chainedRPC)
		}
		return chainedRPC(ctx, req)
	}

	f := func(c context.Context, r proto.Message) (proto.Message, error) {
		return w.pingServer.Ping(ctx, r.(*pb.PingRequest))
	}

	ret, err := chained(ctx, req, f)

	fmt.Println(ret, err)
}

type PingAPIServer struct{}

func (p *PingAPIServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	fmt.Println("CALLED PING")
	return &pb.PingResponse{
		Msg: fmt.Sprintf("PONG: %s", req.GetMsg()),
	}, nil
}
