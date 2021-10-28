package main

import (
	"context"
	"log"
	"net"

	"github.com/antonPalmFolkmann/DISYS_MiniProject2.git/ChatService"

	"google.golang.org/grpc"
)

type Server struct {
	ChatService.UnimplementedChittyChatServiceINServer
	participants []string
	timestamp    int32
	cOUT         ChatService.ChittyChatServiceOUTClient
}

func (s *Server) Publish(ctx context.Context, in *ChatService.PublishMessageRequest) (*ChatService.PublishMessageReply, error) {
	s.timestamp = s.GetLamportTime(in.LamportTime)
	log.Printf("Received publish request, lamport time: %v", s.timestamp)

	log.Printf("Attempt publish, lamport time: %v", s.IncreaseLamportTime()) ////increase, because an event happens
	//TODO: publish code
	s.SendBroadCastRequest(s.cOUT, in.Message, in.ParticipantID)

	var lamporttime = s.IncreaseLamportTime()
	log.Printf("Publish reply, lamport time: %v", lamporttime)
	return &ChatService.PublishMessageReply{
		Reply:       "Publish message reply",
		LamportTime: lamporttime,
	}, nil
}

func (s *Server) Join(ctx context.Context, in *ChatService.JoinRequest) (*ChatService.JoinReply, error) {
	s.timestamp = s.GetLamportTime(in.LamportTime)
	log.Printf("Received join request, lamport time: %v", s.timestamp)

	log.Printf("Attempt join, lamport time: %v", s.IncreaseLamportTime()) //increase, because an event happens
	if contains(s.participants, in.ParticipantID) {
		var lamporttime = s.IncreaseLamportTime()
		log.Printf("Already joined reply, lamport time: %v", lamporttime)
		return &ChatService.JoinReply{
			Reply:       "You are already connected",
			LamportTime: lamporttime,
		}, nil
	} else {
		s.participants = append(s.participants, in.ParticipantID)
		var lamporttime = s.IncreaseLamportTime()
		log.Printf("Join reply, lamport time: %v", lamporttime)
		return &ChatService.JoinReply{
			Reply:       "You are now connected",
			LamportTime: lamporttime,
		}, nil
	}

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (s *Server) Leave(ctx context.Context, in *ChatService.LeaveRequest) (*ChatService.LeaveReply, error) {
	s.timestamp = s.GetLamportTime(in.LamportTime)
	log.Printf("Received leave request, lamport time: %v", s.timestamp)

	log.Printf("Attempt leave, lamport time: %v", s.IncreaseLamportTime()) //increase, because an event happens
	if contains(s.participants, in.ParticipantID) {
		removeByID(s.participants, in.ParticipantID)
		var lamporttime = s.IncreaseLamportTime()
		log.Printf("Leave reply, lamport time: %v", lamporttime)
		return &ChatService.LeaveReply{
			Reply:       "User left the server",
			LamportTime: lamporttime,
		}, nil
	} else {
		var lamporttime = s.IncreaseLamportTime()
		log.Printf("Unknown leave reply, lamport time: %v", lamporttime)
		return &ChatService.LeaveReply{
			Reply:       "Unknown user tried to leave the server",
			LamportTime: lamporttime,
		}, nil
	}

}

func removeByID(participants []string, ID string) []string {
	for i, p := range participants {
		if p == ID {
			return remove(participants, i)
		}
	}
	return participants
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

/*func (s *Server) BroadCast(ctx context.Context, in *ChatService.BroadCastRequest) (*ChatService.BroadCastReply, error) {
	fmt.Printf("Received broadcastrequest request")

	return &ChatService.BroadCastReply{Reply: "BroadCast succeeded"}, nil
}*/

func main() {

	//OUT
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var connOUT *grpc.ClientConn
	connOUT, errOUT := grpc.Dial(":9080", grpc.WithInsecure())
	if errOUT != nil {
		log.Fatalf("Could not connect: %s", errOUT)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer connOUT.Close()

	//  Create new Client from generated gRPC code from proto

	//IN
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	ChatService.RegisterChittyChatServiceINServer(grpcServer, &Server{cOUT: ChatService.NewChittyChatServiceOUTClient(connOUT)})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}

func (s *Server) GetLamportTime(time int32) int32 {
	if s.timestamp > time {
		return s.timestamp + 1
	} else {
		return time + 1
	}
}

func (s *Server) IncreaseLamportTime() int32 {
	s.timestamp = s.timestamp + 1
	return s.timestamp
}

func (s *Server) SendBroadCastRequest(c ChatService.ChittyChatServiceOUTClient, publishMessage string, participantID string) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	var lamportTime = s.IncreaseLamportTime()
	log.Printf("BroadCast Request sent, lamport time: %v", lamportTime)
	message := ChatService.BroadCastRequest{
		Message:       publishMessage,
		LamportTime:   lamportTime,
		ParticipantID: participantID,
	}

	response, err := c.BroadCast(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling BroadCast: %s", err)
	}

	log.Printf("BroadCast response from the Client: %s, lamport time: %v \n", response.Reply, s.GetLamportTime(response.LamportTime))
}
