package auth

import (
	"context"
	"log"

	"quest-manager/internal/core/ports"

	authv1 "github.com/Vi-72/quest-auth/api/grpc/sdk/go/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Factory struct {
	Addr   string
	Client ports.AuthClient
}

func (f *Factory) Create(ctx context.Context) (ports.AuthClient, *grpc.ClientConn, error) {
	if f.Client != nil {
		return f.Client, nil, nil // mock
	}

	if f.Addr == "" {
		log.Println("Auth gRPC address not provided")
		return nil, nil, nil
	}

	conn, err := grpc.NewClient(f.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := authv1.NewAuthServiceClient(conn)
	return NewUserAuthClient(client), conn, nil
}
