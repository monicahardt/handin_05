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
	ackno := &proto.Ack{Ack: success}
	for _, s := range f.servers {

		fmt.Printf("bid : %v\n", bid)
		ack, _ := s.Bid(ctx, bid)
		ackno = ack
		fmt.Printf("Ack: %v\n", ack)
		//What to do here about the bid? Give acknowlegement back
	}
	log.Printf("Client %v recieved acknowlegement %v", bid.Id, ackno.Ack)
	return &proto.Ack{Ack: ackno.Ack}, nil
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
