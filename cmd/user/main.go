package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Tanakaryuki/go-grpc/api/user"
	"github.com/Tanakaryuki/go-grpc/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		port       = flag.Int("port", 50051, "The server port")
		dbHost     = flag.String("db_host", "postgres", "Database host")
		dbPort     = flag.Int("db_port", 5432, "Database port")
		dbUser     = flag.String("db_user", "app_user", "Database user")
		dbPassword = flag.String("db_password", "password", "Database password")
		dbName     = flag.String("db_name", "userdb", "Database name")
	)
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s, err := user.NewServer(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, s)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Printf("User Service is running on port :%d", *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
