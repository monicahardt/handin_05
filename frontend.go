package main

import (
	proto "HANDIN_05/proto"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

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
	//setting the log file
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	flag.Parse()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	frontend := &Frontend{
		name:    "frontdend",
		port:    *port,
		servers: make([]proto.AuctionClient, 0),
		bids:    make([]*proto.Ack, 0),
	}

	go startFrontend(frontend)

	for i := 0; i < 3; i++ {

		conn, err := grpc.Dial("localhost:"+strconv.Itoa(5001+i), grpc.WithTransportCredentials(insecure.NewCredentials()))
		log.Printf("Frontend connected to server at port: %v\n", 5001+i)
		frontend.servers = append(frontend.servers, proto.NewAuctionClient(conn))
		if err != nil {
			log.Printf("Could not connect: %s", err)
		}
		defer conn.Close()

	}
	time.Sleep(time.Second * 20)

	am, _ := frontend.Result(context.Background(), &proto.Empty{})
	log.Printf("Auction closed client %v won with bid %v", am.Id, am.Amount)
	fmt.Printf("Auction closed client %v won with bid %v", am.Id, am.Amount)
	for {

	}
	//Note til eksamen!!!!!!
	//ikke brug port 5000
}

func (f *Frontend) Bid(ctx context.Context, bid *proto.Amount) (*proto.Ack, error) {
	log.Printf("Client %v bid %v", bid.Id, bid.Amount)
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
	log.Printf("Client asked for result")
	highestBid := int32(0)
	highestBidID := int32(0)
	for _, s := range f.servers {

		amount, _ := s.Result(ctx, in)

		if int32(amount.Amount) > highestBid {
			highestBid = amount.Amount
			highestBidID = amount.Id
		}

	}
	log.Printf("client with id %v gave highest bid %v \n", highestBidID, highestBid)
	return &proto.Amount{Amount: highestBid, Id: highestBidID}, nil
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

	}
}

// our enum types
type ack string

const (
	fail      string = "fail"
	success   string = "success"
	exception string = "exception"
)
