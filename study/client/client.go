package main

import (
	"context"
	"io"
	"log"
	"time"

	// "io"

	pb "github.com/PhoenixXiang/grpc-study/study/route"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:10110"
	defaultName = "world"
)

func getOneInfo(c pb.RouteClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetOneInfo(ctx, &pb.Token{Id: 1})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("GetOneInfo: %v", r)

}

func getAllInfo(c pb.RouteClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.GetAllInfo(ctx, &pb.Tokens{Token: []*pb.Token{&pb.Token{Id: 1}, &pb.Token{Id: 2}, &pb.Token{Id: 4}}})
	if err != nil {
		return err
	}
	for {
		p, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetAllInfo(_) = _, %v", c, err)
			return err
		}
		log.Println(p)
	}
	return nil
}

func getSomeInfo(c pb.RouteClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.GetSomeInfo(ctx)
	if err != nil {
		return err
	}
	list := map[int]pb.Token{
		2: {Id: 2},
		4: {Id: 4},
	}
	for _, t := range list {
		stream.Send(&t)
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.getSomeInfo() got error %v, want %v", stream, err, nil)
		return err
	}
	log.Println(reply)
	return nil
}

func getInfo(c pb.RouteClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.GetInfo(ctx)
	if err != nil {
		return err
	}
	list := map[int]pb.Token{
		1: {Id: 1},
		2: {Id: 2},
		4: {Id: 4},
		3: {Id: 3},
		5: {Id: 5},
	}

	waitc := make(chan int)
	go func() {

		for {
			p, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("%v.GetAllInfo(_) = _, %v", c, err)
			}
			log.Printf("getInfo Received: %v", p)
		}
	}()

	for _, t := range list {
		if err := stream.Send(&t); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
			return err
		} else {
			log.Printf("getInfo Send: %v", t)
		}
	}
	stream.CloseSend()
	<-waitc
	return nil
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRouteClient(conn)

	// getOneInfo(c)

	// getAllInfo(c)

	// getSomeInfo(c)

	getInfo(c)

}
