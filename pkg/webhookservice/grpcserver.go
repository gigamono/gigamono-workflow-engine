package webhookservice

import (
	"context"
	"fmt"
	"net"

	"github.com/gigamono/gigamono/pkg/services/proto/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (service *WebhookService) grpcServe(listener net.Listener) error {
	grpcServer := grpc.NewServer() // Create a gRPC service.

	// Register gRPC service.
	generated.RegisterWorkflowWebhookServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	return grpcServer.Serve(listener) // Listen for requests.
}

// SayHello replies with message.
func (service *WebhookService) SayHello(ctx context.Context, msg *generated.Message) (*generated.Message, error) {
	engineMsg := "Webhook Service replies: " + msg.Content
	fmt.Println(engineMsg)
	response := generated.Message{
		Content: engineMsg,
	}
	return &response, nil
}