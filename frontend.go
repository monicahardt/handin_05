package main

import (
	proto "HANDIN_05/proto"
	"context"
	"errors"
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
	auctionIsClosed bool
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
		bids: make([]*proto.Ack, 0),
		auctionIsClosed: false,
	}

	go startFrontend(frontend)

	for i := 0; i < 3; i++ {
		conn, err := grpc.Dial("localhost:"+strconv.Itoa(5000+i), grpc.WithTransportCredentials(insecure.NewCredentials()))
		log.Printf("Frontend connected to server at port: %v\n", 5000+i)
		frontend.servers = append(frontend.servers, proto.NewAuctionClient(conn))
		if err != nil {
			log.Printf("Could not connect: %s", err)
		}
		defer conn.Close()

	}
		time.Sleep(time.Second * 50)
		frontend.auctionIsClosed = true
		am, _ := frontend.Result(context.Background(), &proto.Empty{})
		log.Printf("Auction closed client %v won with bid %v", am.Id, am.Amount)
		fmt.Printf("Auction closed client %v won with bid %v", am.Id, am.Amount)
	
		for{

		}
	//Note til eksamen!!!!!!
	//ikke brug port 5000
}

func (f *Frontend) Bid(ctx context.Context, bid *proto.Amount) (*proto.Ack, error) {
	if f.auctionIsClosed {
		return &proto.Ack{Ack: fail}, errors.New("Auction is closed")
	}
	f.bids = make([]*proto.Ack, 0)
	log.Printf("Client %v bid %v", bid.Id, bid.Amount)
	fmt.Printf("printing number of servers in the slice: %v\n", len(f.servers))
	for index, s := range f.servers {
		ack, err := s.Bid(ctx, bid)
		//tried calling bid on a server that has crashed
		if(err != nil){
			fmt.Println("A server crashed, try to remove from slice")
			f.servers = append(f.servers[:index], f.servers[index+1:]...)
			fmt.Printf("printing size of slice after removal: %v \n", len(f.servers))
		} else {
			f.bids = append(f.bids, ack)
		}

		// hvis ack allerede findes, tæl op
		// hvis ack findes findes i mappet, så tæl dens value en op
	}


	var sCount = 0
	var fCount = 0
	var eCount = 0

	for i := 0; i < len(f.servers); i++ {
		if (f.bids[i].Ack == success) {
			sCount++
		}
		if (f.bids[i].Ack == fail) {
			fCount++
		}
		if (f.bids[i].Ack == exception) {
			eCount++
		}
	}

	fmt.Printf("printing number of success: %v \n", sCount)


	if sCount > (len(f.servers)/2) && sCount != 0 {
		fmt.Println("succes was bigger than half")
		for i := 0; i < len(f.servers); i++ {
			if (f.bids[i].Ack != success) {
				// disconnect the server on f.servers[i]
				f.servers = append(f.servers[:i], f.servers[i+1:]...)
			}
		}
		return &proto.Ack{Ack: success}, nil
	}

	if fCount > (len(f.servers)/2) && fCount != 0 {
		fmt.Println("fails was bigger")
		for i := 0; i < len(f.servers); i++ {
			if (f.bids[i].Ack != fail) {
				// disconnect the server on f.servers[i]
				f.servers = append(f.servers[:i], f.servers[i+1:]...)
			}
		}
		return &proto.Ack{Ack: fail}, nil
	}

	if eCount > (len(f.servers)/2) && eCount != 0 {
		fmt.Println("exception was bigger")
		for i := 0; i < len(f.servers); i++ {
			if (f.bids[i].Ack != exception) {
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
	log.Println("Client asked for result")
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
	fmt.Printf("client with id %v gave highest bid %v \n", highestBidID, highestBid)
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
