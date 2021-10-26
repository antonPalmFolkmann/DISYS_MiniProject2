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
	ChatService.UnimplementedChittyChatServiceServer
}

func (s *Server) BroadCast(ctx context.Context, in *ChatService.BroadCastRequest) (*ChatService.BroadCastReply, error) {
	fmt.Printf("Received broadcastrequest request")
	return &ChatService.BroadCastReply{message: ChatService.BroadCastRequest.message}, nil
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	ChatService.RegisterChittyChatServiceServer(grpcServer, &Server{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func ReceivePublishMessage() {

}
