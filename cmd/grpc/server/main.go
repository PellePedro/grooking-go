package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/user/module/protos"
	"google.golang.org/grpc"
)

const (
	GRPC_PORT = 50051
)

type Server struct {
}

// GRPC

// Unary
func (s *Server) ExecUnary(ctx context.Context, req *protos.UnaryRequest) (*protos.UnaryResponse, error) {
	fmt.Printf("Server Received ExecUnary() with [%s]\n", req.GetMessage())
	return &protos.UnaryResponse{Status: "Server Hello"}, nil
}

// Server Streaming
func (s *Server) ExecServerStreaming(req *protos.ServerStreamingRequest, stream protos.RoutingService_ExecServerStreamingServer) error {

	ch := make(chan interface{}, 100)
	eventChannel.RegisterListener(req.GetClientId(), ch)
	defer func() {
		eventChannel.UnregisterListener(req.GetClientId())
		close(ch)
	}()

	for {
		select {
		case data := <-ch:
			link, ok := data.(Metric)
			if ok {
				res := protos.ServerStreamingResponse{
					Src:     link.src,
					Dst:     link.dst,
					Metric1: int32(link.metric),
				}
				if err := stream.Send(&res); err != nil {
					fmt.Printf("Stream error [%#v]\n", err)
					return nil
				}
				fmt.Printf("Server Streams [%#v]\n", link)
			} else {
				fmt.Println("Failed to decode Metric from event channel")
			}
		}
	}
}

// Client Streaming
func (s *Server) ExecClientStreaming(stream protos.RoutingService_ExecClientStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("EOF while reading client stream: %v", err)
			// finished reading the client stream
			return stream.SendAndClose(&protos.ClientStreamingResponse{})
		}
		if err != nil {
			fmt.Printf("Error while reading client stream: %v", err)
		}
		fmt.Printf("Received Streaming Message message [%#v]\n", req)
	}
}

// BiDi Streaming
func (wu *Server) ExecBiDiStreaming(stream protos.RoutingService_ExecBiDiStreamingServer) error {
	return nil
}

func main() {

	var wg sync.WaitGroup
	wg.Add(1)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", GRPC_PORT))
	if err != nil {
		fmt.Printf("Failed to listen to grpc port: %v", err)
	}

	s := grpc.NewServer()
	protos.RegisterRoutingServiceServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}

	wg.Wait()
}
