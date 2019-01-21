package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/PhoenixXiang/grpc-study/study/route"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":10110"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

var phone1 = pb.Person_PhoneNumber{
	Number: "12345",
	Type:   pb.Person_HOME,
}
var phone2 = pb.Person_PhoneNumber{
	Number: "123456",
	Type:   pb.Person_MOBILE,
}
var phone3 = pb.Person_PhoneNumber{
	Number: "1234567",
	Type:   pb.Person_WORK,
}
var list = map[int]pb.Person{
	1: pb.Person{
		Name:  "xx1",
		Id:    1,
		Email: "xx1@qq.com",
		Phone: []*pb.Person_PhoneNumber{},
		Car:   false,
		Money: 1234.56,
	},
	2: pb.Person{
		Name:  "xx2",
		Id:    2,
		Email: "xx2@qq.com",
		Phone: []*pb.Person_PhoneNumber{&phone1},
		Car:   true,
		Money: 1234.5678,
	},
	3: pb.Person{
		Name:  "xx3",
		Id:    3,
		Email: "xx3@qq.com",
		Phone: []*pb.Person_PhoneNumber{&phone2, &phone3},
		Car:   false,
		Money: 1234.56789,
	},
}

func (s *server) GetOneInfo(ctx context.Context, in *pb.Token) (*pb.Person, error) {
	log.Printf("GetOneInfo Received: %v", in)
	p := list[int(in.Id)]
	return &p, nil
}

func (s *server) GetAllInfo(in *pb.Tokens, stream pb.Route_GetAllInfoServer) error {
	for _, t := range in.Token {
		if p, ok := list[int(t.Id)]; ok {
			if err := stream.Send(&p); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *server) GetSomeInfo(stream pb.Route_GetSomeInfoServer) error {
	startTime := time.Now()
	persons := pb.Persons{Person: []*pb.Person{}}
	for {
		token, err := stream.Recv()

		if err == io.EOF {
			endTime := time.Now()
			log.Printf("GetSomeInfo 耗时: %v s", endTime.Sub(startTime).Seconds())
			return stream.SendAndClose(&persons)
		}

		if err != nil {
			return err
		}

		log.Printf("GetSomeInfo Received: %v", token)
		if p, ok := list[int(token.Id)]; ok {
			persons.Person = append(persons.Person, &p)
		}

	}
}

func (s *server) GetInfo(stream pb.Route_GetInfoServer) error {
	for {
		token, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("GetInfo Received: %v", token)
		if p, ok := list[int(token.Id)]; ok {
			if err := stream.Send(&p); err != nil {
				return err
			}
			log.Printf("GetInfo Send: %v", p)
		}
		
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRouteServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
