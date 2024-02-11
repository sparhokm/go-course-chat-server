package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/sparhokm/go-course-ms-chat-server/pkg/chat_v1"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateIn) (*desc.CreateOut, error) {
	_ = ctx
	log.Printf("Create chat %+v", req)

	return &desc.CreateOut{Id: 1}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteIn) (*emptypb.Empty, error) {
	_ = ctx
	log.Printf("Delete chat: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageIn) (*emptypb.Empty, error) {
	_ = ctx
	log.Printf("Send message %+v", req)

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
