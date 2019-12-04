package main

import (
	"context"
	"fmt"

	"github.com/nametake/chain-interceptor/pb"
)

func main() {
	pingServer := &PingAPIServer{}
	wrapStruct := &WrapStruct{
		pingServer: pingServer,
	}

	wrapStruct.Call()
}

type WrapStruct struct {
	pingServer pb.PingAPIServer
}

func (w *WrapStruct) Call() {
	ctx := context.Background()
	req := &pb.PingRequest{
		Msg: "PING",
	}

	ret, err := w.pingServer.Ping(ctx, req)

	fmt.Println(ret, err)
}

type PingAPIServer struct{}

func (p *PingAPIServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Msg: fmt.Sprintf("PONG: %s", req.GetMsg()),
	}, nil
}
