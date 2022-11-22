package main

import (
	proto "HANDIN_05/proto"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedAuctionServer
	name          string
	port          int
	highestBid    int32
	highestBidder int32
}

var port = flag.Int("port", 0, "server port number") // create the port that recieves the port that the client wants to access to

func main() {
	flag.Parse()

	// highestBid := 0
	// highestBidder := -1

	server := &Server{
		name: "serverName",
		port: *port,
	}

	go startServer(server)

	for {

	}
}

func startServer(server *Server) {
	grpcServer := grpc.NewServer()                                           // create a new grpc server
	listen, err := net.Listen("tcp", "localhost:"+strconv.Itoa(server.port)) // creates the listener

	if err != nil {
		log.Fatalln("Could not start listener")
	}

	log.Printf("Server started")

	proto.RegisterAuctionServer(grpcServer, server)
	serverError := grpcServer.Serve(listen)

	if serverError != nil {
		log.Printf("Could not register server")
	}

}

func (s *Server) Bid(ctx context.Context, bid *proto.Amount) (*proto.Ack, error) {

	fmt.Println("Bid method in server.go was called")
	// tager biddet ind
	// checker om biddet er skarpt større end det registrerede bid
	// hvis det er, returnerer success
	// hvis biddet er mindre end eller lig det registrerede bid
	// returner fail
	// hvis programmet crasher
	// returner exception

	if bid.Amount > s.highestBid {
		s.highestBid = bid.Amount
		s.highestBidder = bid.Id

		return &proto.Ack{Ack: success}, nil
	} else if bid.Amount <= s.highestBid {
		return &proto.Ack{Ack: fail}, nil
	}

	// how do we handle the exception for system crash?

	return &proto.Ack{Ack: exception}, nil
}

func (s *Server) Result(ctx context.Context, in *proto.Empty) (*proto.Amount, error) {
	fmt.Println("result method in server.go was called")
	return &proto.Amount{Amount: s.highestBid, Id: s.highestBidder}, nil
}

// our enum types
type ack string

const (
	fail      string = "fail"
	success   string = "success"
	exception string = "exception"
)
