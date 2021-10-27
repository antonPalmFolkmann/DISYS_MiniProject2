package main

import (
	"context"
	"fmt"

	"log"

	"github.com/antonPalmFolkmann/DISYS_MiniProject2.git/ChatService"

	"google.golang.org/grpc"
)

type Participant struct {
	ID        string
	timestamp int32
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
	c := ChatService.NewChittyChatServiceINClient(conn)

	var p Participant
	p.timestamp = 0
	p.ID = "0"

	//var textmessage = "HEJHEJ"

	//p.SendPublishRequest(c, textmessage)
	p.SendJoinRequest(c)
	p.SendLeaveRequest(c)
}

func (p *Participant) SendPublishRequest(c ChatService.ChittyChatServiceINClient, textmessage string) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	var lamportTime = p.IncreaseLamportTime()
	message := ChatService.PublishMessageRequest{
		LamportTime: lamportTime,
		Message:     textmessage,
	}
	response, err := c.Publish(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Publish: %s", err)
	}

	fmt.Printf("Response from the Server: %s, lamport time: %v \n", response.Reply, p.GetLamportTime(response.LamportTime))
}

func (p *Participant) SendJoinRequest(c ChatService.ChittyChatServiceINClient) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	var lamportTime = p.IncreaseLamportTime()
	log.Printf("Send join request, lamport time: %v", lamportTime)
	message := ChatService.JoinRequest{
		ParticipantID: p.ID,
		LamportTime:   lamportTime,
	}

	response, err := c.Join(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Publish: %s", err)
	}

	log.Printf("Join response from the Server: %s, lamport time: %v \n", response.Reply, p.GetLamportTime(response.LamportTime))
}

func (p *Participant) SendLeaveRequest(c ChatService.ChittyChatServiceINClient) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	var lamportTime = p.IncreaseLamportTime()
	message := ChatService.LeaveRequest{
		ParticipantID: p.ID,
		LamportTime:   lamportTime,
	}

	response, err := c.Leave(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Publish: %s", err)
	}

	log.Printf("Leave response from the Server: %s, lamport time: %v \n", response.Reply, p.GetLamportTime(response.LamportTime))
}

func (p *Participant) GetLamportTime(time int32) int32 {
	if p.timestamp > time {
		return p.timestamp + 1
	} else {
		return time + 1
	}
}

func (p *Participant) IncreaseLamportTime() int32 {
	p.timestamp = p.timestamp + 1
	return p.timestamp
}
