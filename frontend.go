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
	"google.golang.org/grpc/credentials/insecure"
)

type Frontend struct {
	proto.UnimplementedAuctionServer
	name    string
	port    int
	servers []proto.AuctionClient
}

var port = flag.Int("port", 0, "server port number") // create the port that recieves the port that the client wants to access to

func main() {
	flag.Parse()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	frontend := &Frontend{
		name:    "frondend",
		port:    *port,
		servers: make([]proto.AuctionClient, 0),
	}
	go startFrontend(frontend)
	log.Printf("Frontend started")
	fmt.Printf("Frontend started at port: %v", port)

	fmt.Printf("We get here!")
	for i := 0; i < 3; i++ {

		fmt.Printf("Trying to dial: %v\n", 5001+i)
		conn, err := grpc.Dial("localhost:"+strconv.Itoa(5001+i), grpc.WithTransportCredentials(insecure.NewCredentials()))
		fmt.Printf("err : %v", err)
		// serverClient := proto.NewAuctionClient(conn)
		fmt.Printf("connection made to %v\n", 5001+i)
		frontend.servers = append(frontend.servers, proto.NewAuctionClient(conn))
		if err != nil {
			fmt.Printf("Could not connect: %s", err)
		}
		defer conn.Close()

	}

	for {

	}
	//Note til eksamen!!!!!!
	//ikke brug port 5000
}

func (f *Frontend) Bid(ctx context.Context, bid *proto.Amount) (*proto.Ack, error) {
	fmt.Printf("Called method bid")
	fmt.Printf("slice length: %v\n", len(f.servers))
	for _, s := range f.servers {
		fmt.Println("There was a server in the slice")
		fmt.Printf("bid : %v\n", bid)
		ack, _ := s.Bid(ctx, bid)
		fmt.Printf("Ack: %v\n", ack)
		//What to do here about the bid? Give acknowlegement back
	}
	fmt.Println("Returns")
	return &proto.Ack{Ack: success}, nil
}

func (f *Frontend) Result(ctx context.Context, in *proto.Empty) (*proto.Amount, error) {
	fmt.Printf("Called method result")
	return &proto.Amount{Amount: 5, Id: 4}, nil
}

func startFrontend(frontend *Frontend) {
	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(frontend.port))

	if err != nil {
		log.Fatalf("Could not create the frontend %v", err)
	}
	log.Printf("Started frontend at port: %d\n", frontend.port)

	// Register the grpc server and serve its listener
	proto.RegisterAuctionServer(grpcServer, frontend)

	serveError := grpcServer.Serve(listener)
	fmt.Printf("nedern")
	if serveError != nil {
		log.Fatalf("Could not serve listener frontend")
		fmt.Printf("Could not serve listener frontend")

	}
}

// our enum types
type ack string

const (
	fail      string = "fail"
	success   string = "success"
	exception string = "exception"
)
