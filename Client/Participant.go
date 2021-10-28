package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"log"

	"github.com/antonPalmFolkmann/DISYS_MiniProject2.git/ChatService"

	"google.golang.org/grpc"
)

type Participant struct {
	ID              string
	timestamp       int32
	connectToServer bool
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

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please write your ID")
	input, _ := reader.ReadString('\n')
	// convert CRLF to LF
	input = strings.Replace(input, "\n", "", -1)

	var p Participant
	p.timestamp = 0
	p.ID = input
	p.connectToServer = true
	p.SendJoinRequest(c)

	for {

		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		if strings.Compare("/leave", input) == 0 {
			p.SendLeaveRequest(c)
		} else {
			p.SendPublishRequest(c, input)
		}

		if !p.connectToServer {
			break
		}
	}

	//var textmessage = "HEJHEJ"

	//p.SendPublishRequest(c, textmessage)

}

func (p *Participant) SendPublishRequest(c ChatService.ChittyChatServiceINClient, textmessage string) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	var lamportTime = p.IncreaseLamportTime()
	log.Printf("Publish Request sent, lamport time: %v", lamportTime)
	message := ChatService.PublishMessageRequest{
		LamportTime: lamportTime,
		Message:     textmessage,
	}
	response, err := c.Publish(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Publish: %s", err)
	}

	log.Printf("Response from the Server: %s, lamport time: %v \n", response.Reply, p.GetLamportTime(response.LamportTime))
}

func (p *Participant) SendJoinRequest(c ChatService.ChittyChatServiceINClient) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	var lamportTime = p.IncreaseLamportTime()
	log.Printf("Join Request sent, lamport time: %v", lamportTime)
	message := ChatService.JoinRequest{
		ParticipantID: p.ID,
		LamportTime:   lamportTime,
	}

	response, err := c.Join(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Join: %s", err)
	}

	log.Printf("Join response from the Server: %s, lamport time: %v \n", response.Reply, p.GetLamportTime(response.LamportTime))
}

func (p *Participant) SendLeaveRequest(c ChatService.ChittyChatServiceINClient) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	var lamportTime = p.IncreaseLamportTime()
	log.Printf("Leave Request sent, lamport time: %v", lamportTime)
	message := ChatService.LeaveRequest{
		ParticipantID: p.ID,
		LamportTime:   lamportTime,
	}

	response, err := c.Leave(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Leave: %s", err)
	}

	log.Printf("Leave response from the Server: %s, lamport time: %v \n", response.Reply, p.GetLamportTime(response.LamportTime))
	p.connectToServer = false
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
