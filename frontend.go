package main

import(
	proto "HANDIN_05/proto"
)


var server[] *proto.AuctionServer
var client *proto.AuctionClient



func main(){

	for server := range server{
		server.Bid(client.ClientConn)
	}
	
	func (s *Server) Bid(ctx context.Context, bid *proto.Amount) (*proto.Ack, error) {
	
	}

}

