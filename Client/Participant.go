package main

import (
	"context"
	"fmt"

	"log"

	"github.com/antonPalmFolkmann/DISYS_MiniProject2.git/ChatService"

	"google.golang.org/grpc"
)

type Participant struct {
	
}

func main() {
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := ChatService.ChittyChatServiceClient(conn)

	SendPublishRequest(c)
}

func SendPublishRequest(c ChatService.ChittyChatServiceClient) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	message := ChatService.PublishMessageRequest{}

	response, err := c.Publish(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Publish: %s", err)
	}

	fmt.Printf("Response from the Server: %s \n", response.Reply)
}

func (p *Participant) ReadBroadCastChannel(){

}
