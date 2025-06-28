package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/Git5737/lexchanger/proto/chat/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiseServer
	clients map[string]pb.ChatServise_EventStreamServer
}

func newServer() *server {
	return &server{clients: make(map[string]pb.ChatServise_EventStreamServer)}
}

func (s *server) EventStream(stream pb.ChatServise_EventStreamServer) error {
	var clientName string

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("receive error: %v", err)
			return err
		}

		switch ev := in.Event.(type) {
		case *pb.Events_ClientLogin:
			clientName = ev.ClientLogin.Name
			s.clients[clientName] = stream
			log.Printf("User %s joined", clientName)

		case *pb.Events_ClientMessage:
			msg := fmt.Sprintf("%s: %s", ev.ClientMessage.Name, ev.ClientMessage.Message)
			log.Println(msg)

			for _, client := range s.clients {
				//if name != clientName {
				client.Send(&pb.Events{
					Event: &pb.Events_ClientMessage{
						ClientMessage: &pb.Events_Message{
							Name:    ev.ClientMessage.Name,
							Message: ev.ClientMessage.Message,
						},
					},
				})
				//}
			}

		case *pb.Events_ClientLogout:
			delete(s.clients, ev.ClientLogout.Name)
			log.Printf("User %s left", ev.ClientLogout.Name)
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiseServer(grpcServer, newServer())

	log.Println("Server is running at :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
