package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/antonPalmFolkmann/DISYS_MiniProject2.git/ChatService"

	"google.golang.org/grpc"
)

type Participant struct {
	ID              string
	timestamp       int32
	connectToServer bool
	ChatService.UnimplementedChittyChatServiceOUTServer
}

func main() {
	file, err := os.OpenFile("../info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	//IN
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please write your port")
	input, _ := reader.ReadString('\n')
	// convert CRLF to LF
	input = strings.Replace(input, "\n", "", -1)
	input = strings.Replace(input, "\r", "", -1)

	var p Participant
	p.timestamp = 0
	p.ID = input
	p.connectToServer = true

	portString := ":" + input
	// Create listener tcp on port *input*
	go p.BroadCastListen(portString, p.ID)

	//for publish and join and leave
	var conn *grpc.ClientConn
	conn, errIN := grpc.Dial(":7080", grpc.WithInsecure())
	if errIN != nil {
		log.Fatalf("Could not connect: %s", errIN)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := ChatService.NewChittyChatServiceINClient(conn)

	p.SendJoinRequest(c)

	for {

		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)
		input = strings.Replace(input, "\r", "", -1)

		if strings.Compare("/leave", input) == 0 {
			p.SendLeaveRequest(c)
		} else if len(input) <= 128 {
			p.SendPublishRequest(c, input, true, false, false)
		} else {
			fmt.Println("Invalid message")
		}

		if !p.connectToServer {
			break
		}
	}
}

func (p *Participant) BroadCastListen(portString string, pID string) {
	list, err := net.Listen("tcp", portString)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", p.ID, err)
	}
	grpcServer := grpc.NewServer()
	ChatService.RegisterChittyChatServiceOUTServer(grpcServer, &Participant{ID: pID})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func (p *Participant) SendPublishRequest(c ChatService.ChittyChatServiceINClient, textmessage string, isPublish bool, isJoin bool, isLeave bool) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	var lamportTime = p.IncreaseLamportTime()

	if !isLeave {
		log.Printf("Publish Request sent, lamport time: %v", lamportTime)
	}

	message := ChatService.Message{
		LamportTime:   lamportTime,
		Message:       textmessage,
		ParticipantID: p.ID,
		IsPublish:     isPublish,
		IsJoin:        isJoin,
		IsLeave:       isLeave,
	}
	response, err := c.Publish(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Publish: %s", err)
	}

	if !isLeave {
		log.Printf("Response from the Server: %s, lamport time: %v \n", response.Reply, p.GetLamportTime(response.LamportTime))
	}

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

	stringToPublish := "Participant " + p.ID + " joined Chitty-Chat at Lamport time " + strconv.Itoa(int(p.timestamp))

	p.SendPublishRequest(c, stringToPublish, false, true, false)
}

func (p *Participant) SendLeaveRequest(c ChatService.ChittyChatServiceINClient) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	var lamportTime = p.IncreaseLamportTime()
	//log.Printf("Leave Request sent, lamport time: %v", lamportTime)

	message := ChatService.LeaveRequest{
		ParticipantID: p.ID,
		LamportTime:   lamportTime,
	}

	log.Printf("Leave Request sent, lamport time: %v", lamportTime)
	stringToPublish := "Participant " + p.ID + " left Chitty-Chat at Lamport time " + strconv.Itoa(int(p.timestamp))
	p.SendPublishRequest(c, stringToPublish, false, false, true)

	response, err := c.Leave(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Leave: %s", err)
	}

	log.Printf("Leave response from the Server: %s, lamport time: %v \n", response.Reply, p.GetLamportTime(response.LamportTime))
	p.connectToServer = false

}

func (p *Participant) GetLamportTime(time int32) int32 {
	if p.timestamp > time {
		p.timestamp += 1
	} else {
		p.timestamp = time + 1

	}
	return p.timestamp
}

func (p *Participant) IncreaseLamportTime() int32 {
	p.timestamp = p.timestamp + 1
	return p.timestamp
}

func (p *Participant) BroadCast(ctx context.Context, in *ChatService.BroadCastRequest) (*ChatService.BroadCastReply, error) {
	if !in.Message.IsLeave {
		log.Printf("Received broadcast request, lamport time: %v", p.GetLamportTime(in.LamportTime))

		log.Printf("Attempt BroadCast, lamport time: %v", p.IncreaseLamportTime()) //increase, because an event happens
	}

	if in.Message.IsPublish {
		fmt.Printf("%v said: %v, lamport time: %v\n", in.ParticipantID, in.Message.Message, p.timestamp)
		log.Printf("%v said: %v, lamport time: %v", in.ParticipantID, in.Message.Message, p.timestamp)
	}
	if in.Message.IsJoin {
		fmt.Println(in.Message.Message)
		log.Printf(in.Message.Message)
	}
	if in.Message.IsLeave {
		if in.Message.ParticipantID != p.ID {

			fmt.Println(in.Message.Message)
			log.Printf(in.Message.Message)

			var lamporttime = p.IncreaseLamportTime()
			return &ChatService.BroadCastReply{
				Reply:       "BroadCast message reply",
				LamportTime: lamporttime,
			}, nil
		}
	}

	var lamporttime = p.IncreaseLamportTime()
	if !in.Message.IsLeave {
		log.Printf("BroadCast reply, lamport time: %v", lamporttime)
	}
	return &ChatService.BroadCastReply{
		Reply:       "BroadCast message reply",
		LamportTime: lamporttime,
	}, nil

}
