package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/user/module/protos"
	"google.golang.org/grpc"
)

func doUnary(c protos.RoutingServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &protos.UnaryRequest{
		Message: "Client says Hello",
	}
	res, err := c.ExecUnary(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling ExecUnary RPC: %v", err)
	}
	fmt.Printf("Response from Server.doUnary() call: [%s]\n", res.Status)
}

func doClientStreaming(c protos.RoutingServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*protos.ClientStreamingRequest{
		{
			Src: "10.10.0.1",
		},
		{
			Src: "10.10.0.2",
		},
		{
			Src: "10.10.0.3",
		},
	}

	stream, err := c.ExecClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("error while calling DoClientStreaming: %#v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from Server: %v", err)
	}
	fmt.Printf("ClientStreaming Response: %v\n", res)
}

func doServerStreaming(c protos.RoutingServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &protos.ServerStreamingRequest{ClientId: "1"}

	resStream, err := c.ExecServerStreaming(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		_ = msg.GetSrc()
		_ = msg.GetDst()
		_ = msg.GetMetric1()
		_ = msg.GetMetric2()
		_ = msg.GetMetric3()
		fmt.Printf("Received event src[%s] dst[%s] metric[%d]\n", msg.GetSrc(), msg.GetDst(), msg.GetMetric1())
	}
}

func doBiDiStreaming(c protos.RoutingServiceClient) {

	// we create a stream by invoking the client
	stream, err := c.ExecBiDiStreaming(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}
	requests := []*protos.BiDiStreamingRequest{
		{}, {}, {},
	}

	waitc := make(chan struct{})
	// sending streaming messages
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %#v\n", res)
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
	_ = stream
}

//func generateEvent(pipe chan protos.ClientStreamingRequest)

func main() {

	var wg sync.WaitGroup
	wg.Add(1)

	fmt.Println("Hello I'm a client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := protos.NewRoutingServiceClient(cc)

	doUnary(c)
	//doClientStreaming(c)
	doServerStreaming(c)
	// doBiDiStreaming(c)
	// doUnaryWithDeadline(c, 5*time.Second) // should complete
	// doUnaryWithDeadline(c, 1*time.Second) // should timeout
	wg.Wait()
}
