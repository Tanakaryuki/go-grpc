package order

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	pb "github.com/Tanakaryuki/go-grpc/api/order"
	_ "github.com/lib/pq"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	db *sql.DB
}

func NewServer(host string, port int, user, password, dbname string) (*Server, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// 確認のためにデータベースに接続
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// テーブルが存在しない場合は作成
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS orders (
        id SERIAL PRIMARY KEY,
        user_id INTEGER,
        product_id INTEGER,
        quantity INTEGER,
        status VARCHAR(50)
    );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &Server{db: db}, nil
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var id int
	// 初期ステータスを"Pending"とする
	err := s.db.QueryRow("INSERT INTO orders (user_id, product_id, quantity, status) VALUES ($1, $2, $3, $4) RETURNING id",
		req.UserId, req.ProductId, req.Quantity, "Pending").Scan(&id)
	if err != nil {
		log.Printf("Error inserting order: %v", err)
		return nil, err
	}
	return &pb.CreateOrderResponse{Id: fmt.Sprintf("%d", id)}, nil
}

func (s *Server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	var id, userId, productId, quantity int
	var status string
	err := s.db.QueryRow("SELECT id, user_id, product_id, quantity, status FROM orders WHERE id = $1", req.Id).
		Scan(&id, &userId, &productId, &quantity, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		log.Printf("Error fetching order: %v", err)
		return nil, err
	}
	return &pb.GetOrderResponse{
		Id:        fmt.Sprintf("%d", id),
		UserId:    fmt.Sprintf("%d", userId),
		ProductId: fmt.Sprintf("%d", productId),
		Quantity:  int32(quantity),
		Status:    status,
	}, nil
}
