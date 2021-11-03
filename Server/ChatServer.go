package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/antonPalmFolkmann/DISYS_MiniProject2.git/ChatService"

	"google.golang.org/grpc"
)

type Server struct {
	ChatService.UnimplementedChittyChatServiceINServer
	participants []string
	timestamp    int32
	//messages []*ChatService.Message
}

func (s *Server) Publish(ctx context.Context, in *ChatService.Message) (*ChatService.PublishMessageReply, error) {
	s.timestamp = s.GetLamportTime(in.LamportTime)
	log.Printf("Received publish request, lamport time: %v", s.timestamp)

	log.Printf("Attempt publish, lamport time: %v", s.IncreaseLamportTime()) ////increase, because an event happens
	s.SendBroadCastRequest(in)

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
		s.participants = removeByID(s.participants, in.ParticipantID)

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
	newParparticipants := make([]string, 0)

	for _, p := range participants {
		if p != ID {
			newParparticipants = append(newParparticipants, p)
		}
	}

	return newParparticipants
}

func (s *Server) SendBroadCastRequest(textmessage *ChatService.Message) {
	var lamportTime = s.IncreaseLamportTime()
	log.Printf("BroadCast Request sent, lamport time: %v", lamportTime)

	for _, p := range s.participants {
		portString := ":" + p
		// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}

		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()

		//  Create new Client from generated gRPC code from proto
		c := ChatService.NewChittyChatServiceOUTClient(conn)

		var lamportTime = s.IncreaseLamportTime()
		log.Printf("BroadCast Request sent to %v, lamport time: %v", p, lamportTime)
		message := ChatService.BroadCastRequest{
			Message:       textmessage,
			ParticipantID: textmessage.ParticipantID,
			LamportTime:   lamportTime,
		}
		response, err := c.BroadCast(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling BroadCast: %s", err)
		}

		log.Printf("Response from %v: %s, lamport time: %v \n", p, response.Reply, s.GetLamportTime(response.LamportTime))
	}

}

func main() {
	file, err := os.OpenFile("../info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	//IN
	// Create listener tcp on port 7080
	list, err := net.Listen("tcp", ":7080")
	if err != nil {
		log.Fatalf("Failed to listen on port 7080: %v", err)
	}
	grpcServer := grpc.NewServer()
	ChatService.RegisterChittyChatServiceINServer(grpcServer, &Server{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}

func (s *Server) GetLamportTime(time int32) int32 {
	if s.timestamp > time {
		s.timestamp += 1
	} else {
		s.timestamp = time + 1
	}
	return s.timestamp
}

func (s *Server) IncreaseLamportTime() int32 {
	s.timestamp = s.timestamp + 1
	return s.timestamp
}
