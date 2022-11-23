package main

import (
	proto "HANDIN_05/proto"
	"context"
	"flag"
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

	log.Printf("Server started at port %v", server.port)

	proto.RegisterAuctionServer(grpcServer, server)
	serverError := grpcServer.Serve(listen)

	if serverError != nil {
		log.Printf("Could not register server")
	}

}

// the bid method in the server checks if the bid placed by the client is a succesful bid, failed bid or an exception
// updates the highestBid and highestBidder value if the placed bid is a success
// returns the Ack enum corresponding to the bid
func (s *Server) Bid(ctx context.Context, bid *proto.Amount) (*proto.Ack, error) {

	if bid.Amount <= s.highestBid {
		log.Println("return fail")
		return &proto.Ack{Ack: fail}, nil
	} else if bid.Amount > s.highestBid {
		s.highestBid = bid.Amount
		s.highestBidder = bid.Id
		log.Println("New highest bid, return success")
		return &proto.Ack{Ack: success}, nil
	}
	log.Println("return exception")

	return &proto.Ack{Ack: exception}, nil
}

func (s *Server) Result(ctx context.Context, in *proto.Empty) (*proto.Amount, error) {
	log.Println("result method in server was called")
	return &proto.Amount{Amount: s.highestBid, Id: s.highestBidder}, nil
}

// our enum types
type ack string

const (
	fail      string = "fail"
	success   string = "success"
	exception string = "exception"
)
