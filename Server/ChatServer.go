package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/antonPalmFolkmann/DISYS_MiniProject2.git/ChatService"

	"google.golang.org/grpc"
)

type Server struct {
	ChatService.UnimplementedChittyChatServiceINServer
	participants []string
	timestamp    int32
}

func (s *Server) PublishRequest(ctx context.Context, in *ChatService.PublishMessageRequest) (*ChatService.PublishMessageReply, error) {
	log.Printf("Received publish request")
	s.timestamp = s.GetLamportTime(in.LamportTime)
	return &ChatService.PublishMessageReply{
		Reply:       "Your reply here",
		LamportTime: s.IncreaseLamportTime(),
	}, nil
}

func (s *Server) JoinRequest(ctx context.Context, in *ChatService.JoinRequest) (*ChatService.JoinReply, error) {
	log.Printf("Received publish request")
	s.participants = append(s.participants, in.Participant.GetID())
	s.timestamp = s.GetLamportTime(in.LamportTime)
	return &ChatService.JoinReply{
		Reply:       "Your reply here",
		LamportTime: s.IncreaseLamportTime(),
	}, nil
}

func (s *Server) LeaveRequest(ctx context.Context, in *ChatService.LeaveRequest) (*ChatService.LeaveReply, error) {
	log.Printf("Received publish request")
	removeByID(s.participants, in.Participant.ID)
	s.timestamp = s.GetLamportTime(in.LamportTime)
	return &ChatService.LeaveReply{
		Reply:       "Your reply here",
		LamportTime: s.IncreaseLamportTime(),
	}, nil
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

func (s *Server) BroadCast(ctx context.Context, in *ChatService.BroadCastRequest) (*ChatService.BroadCastReply, error) {
	fmt.Printf("Received broadcastrequest request")

	return &ChatService.BroadCastReply{Reply: "BroadCast succeeded"}, nil
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	ChatService.RegisterChittyChatServiceINServer(grpcServer, &Server{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func ReceivePublishMessage() {

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
