package main

import (
	"bufio"
	"context"
	"flag"
	"grpc_kursus/proto"
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

		log.Printf("Client inputted %s\n", input)

		timeMessage, err := serverConnection.GetTime(context.Background(), &proto.AskForTimeMessage{ClientId: int64(client.id)})

		if err != nil {
			log.Printf("Could not get time")
		}

		log.Printf("Server says that the time is %s\n", timeMessage.Time)
	}

}

func getServerConnection() proto.TimeAskServiceClient {

	connection, err := grpc.Dial(":"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials())) // remember to put the last line in the dial function

	if err != nil {
		log.Fatalln("Could not dial")
	}

	log.Printf("Dialed")

	return proto.NewTimeAskServiceClient(connection)
}
