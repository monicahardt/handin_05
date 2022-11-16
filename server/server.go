package main

import (
	"context"
	"flag"
	proto "grpc_kursus/proto"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedTimeAskServiceServer
	name string
	port int
}

var port = flag.Int("port", 8080, "server port number") // create the port that recieves the port that the client wants to access to

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
	grpcServer := grpc.NewServer()                                  // create a new grpc server
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(server.port)) // creates the listener

	if err != nil {
		log.Fatalln("Count not start listener")
	}

	log.Printf("Server started")

	proto.RegisterTimeAskServiceServer(grpcServer, server)
	serverError := grpcServer.Serve(listen)

	if serverError != nil {
		log.Printf("Could not register server")
	}

}

func (c *Server) GetTime(ctx context.Context, in *proto.AskForTimeMessage) (*proto.TimeMessage, error) {
	log.Printf("Client with ID %d asked for the time\n", in.ClientId)

	return &proto.TimeMessage{
		Time:       time.Now().String(),
		ServerName: c.name,
	}, nil
}
