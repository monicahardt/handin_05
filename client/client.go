package main

import (
	"HANDIN_05/proto"
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	id         int
	portNumber int
}

var (
	clientPort   = flag.Int("cPort", 0, "client port number")
	frontendPort = flag.Int("fPort", 0, "frontend port number")
)

func main() {

	flag.Parse()

	client := &Client{
		id:         *clientPort,
		portNumber: *clientPort,
	}

	go connectToFrontend(client)
	log.Printf("Client: %v connected to frontend", client.id)

	for {

	}
}

type Frontend struct {
}

func connectToFrontend(client *Client) {
	FrontendClient := getFrontendConnection()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		//this is the input of the client who want to make a bid or see the result of the auction
		method := scanner.Text()

		if method == "bid" {
			scanner.Scan()
			amountToBid, _ := strconv.ParseInt(scanner.Text(), 10, 0)
			_, err := FrontendClient.Bid(context.Background(), &proto.Amount{Amount: int32(amountToBid), Id: int32(client.id)})

			if(err != nil){
		
				fmt.Printf("Auction has closed, you can ask for the result only")
			}
		} else if method == "result" {
			amount,_ := FrontendClient.Result(context.Background(), &proto.Empty{})
			fmt.Printf("The current highest bidder has id: %v with the highest bid: %v \n", amount.Id, amount.Amount)
		} else {
			fmt.Println("Invalid")
		}

	}
}

func getFrontendConnection() proto.AuctionClient {

	connection, err := grpc.Dial(":"+strconv.Itoa(*frontendPort), grpc.WithTransportCredentials(insecure.NewCredentials())) // remember to put the last line in the dial function

	if err != nil {
		log.Fatalln("Could not dial")
	}

	log.Printf("Dialed")

	return proto.NewAuctionClient(connection)
}
