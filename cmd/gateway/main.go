package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Tanakaryuki/go-grpc/api/inventory"
	"github.com/Tanakaryuki/go-grpc/api/order"
	"github.com/Tanakaryuki/go-grpc/api/product"
	"github.com/Tanakaryuki/go-grpc/api/user"
	gw "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	var (
		grpcUserEndpoint      = flag.String("user_grpc", "user-service:50051", "gRPC User service endpoint")
		grpcProductEndpoint   = flag.String("product_grpc", "product-service:50052", "gRPC Product service endpoint")
		grpcOrderEndpoint     = flag.String("order_grpc", "order-service:50053", "gRPC Order service endpoint")
		grpcInventoryEndpoint = flag.String("inventory_grpc", "inventory-service:50054", "gRPC Inventory service endpoint")
		httpPort              = flag.Int("http_port", 8080, "HTTP port")
	)
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := gw.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Register User Service
	err := user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, *grpcUserEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register UserService handler: %v", err)
	}

	// Register Product Service
	err = product.RegisterProductServiceHandlerFromEndpoint(ctx, mux, *grpcProductEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register ProductService handler: %v", err)
	}

	// Register Order Service
	err = order.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, *grpcOrderEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register OrderService handler: %v", err)
	}

	// Register Inventory Service
	err = inventory.RegisterInventoryServiceHandlerFromEndpoint(ctx, mux, *grpcInventoryEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register InventoryService handler: %v", err)
	}

	log.Printf("API Gateway is running on port :%d", *httpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
