package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Tanakaryuki/go-grpc/api/inventory"
	"github.com/Tanakaryuki/go-grpc/internal/inventory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		port       = flag.Int("port", 50054, "The server port")
		dbHost     = flag.String("db_host", "postgres", "Database host")
		dbPort     = flag.Int("db_port", 5432, "Database port")
		dbUser     = flag.String("db_user", "inventoryuser", "Database user")
		dbPassword = flag.String("db_password", "password", "Database password")
		dbName     = flag.String("db_name", "inventorydb", "Database name")
	)
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s, err := inventory.NewServer(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, s)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Printf("Inventory Service is running on port :%d", *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
