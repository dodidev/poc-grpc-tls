// server/server.go

package main

import (
	"context"
	"flag"
	"fmt"
	"grpc-tls/server/protos"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 50051, "The server port")
)

type server struct {
	protos.UnimplementedDataServiceServer
}

func (s *server) GetJsonData(ctx context.Context, req *protos.Request) (*protos.Response, error) {
	// Prepare the JSON data (just an example)

	jsonData := `[
    {"name": "John Doe", "email": "john.doe@example.com", "phone": "123-456-7890", "address": "123 Main St, Springfield, IL"},
    {"name": "Jane Smith", "email": "jane.smith@example.com", "phone": "987-654-3210", "address": "456 Elm St, Springfield, IL"},
    {"name": "Sam Wilson", "email": "sam.wilson@example.com", "phone": "555-555-5555", "address": "789 Oak St, Springfield, IL"}
  ]`

	// Send back the JSON data as a response
	return &protos.Response{JsonData: jsonData}, nil
}

func main() {
	// Start the gRPC server
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("192.168.1.220:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	protos.RegisterDataServiceServer(grpcServer, &server{})

	fmt.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
