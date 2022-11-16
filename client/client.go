package main

import (
	"HANDIN_05/proto"
	"bufio"
	"context"
	"flag"
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
	clientPort = flag.Int("cPort", 8081, "client port number")
	serverPort = flag.Int("sPort", 8080, "server port number")
)

func main() {

	flag.Parse()

	client := &Client{
		id:         1,
		portNumber: *clientPort,
	}

	go startClient(client)

	for {

	}
}

func startClient(client *Client) {
	serverConnection := getServerConnection()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()

		number := ParseInt(input, 10, 32)

		amount := &proto.Amount{
			Amount: number,
			Id:     int32(client.id),
		}

		bid, err := serverConnection.Bid(context.Background(), &proto.Amount{Amount: number, Id: int32(client.id)})

		log.Printf("Bid returned with %s\n", bid.Ack)

		//timeMessage, err := serverConnection.GetTime(context.Background(), &proto.AskForTimeMessage{ClientId: int64(client.id)})

		if err != nil {
			log.Printf("Could not get time")
		}

		//log.Printf("Server says that the time is %s\n", timeMessage.Time)
	}

}
func makeBid() {}

func getServerConnection() proto.AuctionClient {

	connection, err := grpc.Dial(":"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials())) // remember to put the last line in the dial function

	if err != nil {
		log.Fatalln("Could not dial")
	}

	log.Printf("Dialed")

	return proto.NewAuctionClient(connection)
}
