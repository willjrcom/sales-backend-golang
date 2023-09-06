package grpcimpl

import (
	"context"
	"flag"
	"fmt"

	"github.com/willjrcom/sales-backend-go/internal/infra/pb/orderpb"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	orderpb.UnimplementedOrderServiceServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) GetAllOrder(in *orderpb.BlankMessage, stream orderpb.OrderService_GetAllOrderServer) error {
	stream.Send(&orderpb.Order{Name: "Hello"})
	return nil
}

func (s *Server) Hello(ctx context.Context, in *orderpb.Message) (*orderpb.Message, error) {
	fmt.Println(in.Name)
	return &orderpb.Message{Name: in.GetName()}, nil
}
