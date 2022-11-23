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
	bids    []*proto.Ack
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
		bids:    make([]*proto.Ack, 0),
	}

	go startFrontend(frontend)
	log.Printf("Frontend started")
	fmt.Printf("Frontend started at port: %v", port)

	fmt.Printf("We get here!")
	for i := 0; i < 2; i++ {
		fmt.Println("Yoooooooooooo")
		fmt.Printf("Trying to dial: %v\n", 5000+i)
		conn, err := grpc.Dial("localhost:"+strconv.Itoa(5000+i), grpc.WithTransportCredentials(insecure.NewCredentials()))

		newServer := proto.NewAuctionClient(conn)
		frontend.servers = append(frontend.servers, newServer)
		if err != nil {
			fmt.Printf("Could not connect: %s", err)
		}
		defer conn.Close()

	}

	for {

	}
}

func (f *Frontend) Bid(ctx context.Context, bid *proto.Amount) (*proto.Ack, error) {
	fmt.Printf("Called method bid")
	for _, s := range f.servers {
		fmt.Println("There was a server in the slice")

		ack, _ := s.Bid(ctx, bid)
		// hvis ack allerede findes, tæl op
		// hvis ack findes findes i mappet, så tæl dens value en op

		f.bids = append(f.bids, ack)

		fmt.Println(f.bids) // check if the bids are added to the slice
		//What to do here about the bid? Give acknowlegement back
	}

	var sCount = 0
	var fCount = 0
	var eCount = 0

	for i := 0; i < len(f.servers); i++ {
		if (f.bids[i] == &proto.Ack{Ack: success}) {
			sCount++
		}
		if (f.bids[i] == &proto.Ack{Ack: fail}) {
			fCount++
		}
		if (f.bids[i] == &proto.Ack{Ack: exception}) {
			eCount++
		}
	}

	if sCount > (len(f.servers)/2) && sCount != 0 {
		for i := 0; i < len(f.servers); i++ {
			if (f.bids[i] != &proto.Ack{Ack: success}) {
				// disconnect the server on f.servers[i]
				f.servers = append(f.servers[:i], f.servers[i+1:]...)
			}
		}
		return &proto.Ack{Ack: success}, nil
	}

	if fCount > (len(f.servers)/2) && fCount != 0 {
		for i := 0; i < len(f.servers); i++ {
			if (f.bids[i] != &proto.Ack{Ack: fail}) {
				// disconnect the server on f.servers[i]
				f.servers = append(f.servers[:i], f.servers[i+1:]...)
			}
		}
		return &proto.Ack{Ack: fail}, nil
	}

	if eCount > (len(f.servers)/2) && eCount != 0 {
		for i := 0; i < len(f.servers); i++ {
			if (f.bids[i] != &proto.Ack{Ack: exception}) {
				// disconnect the server on f.servers[i]
				f.servers = append(f.servers[:i], f.servers[i+1:]...)
			}
		}
		return &proto.Ack{Ack: exception}, nil
	}

	// else everyone answered something different and therefore they're all faulty
	fmt.Println("At the end of bid, we're fuckd")
	return &proto.Ack{Ack: exception}, nil
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

type ack string

const (
	fail      string = "fail"
	success   string = "success"
	exception string = "exception"
)
