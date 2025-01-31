package product

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	pb "github.com/Tanakaryuki/go-grpc/api/product"
	_ "github.com/lib/pq"
)

type Server struct {
	pb.UnimplementedProductServiceServer
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
    CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        price FLOAT
    );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &Server{db: db}, nil
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	var id int
	err := s.db.QueryRow("INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id",
		req.Name, req.Price).Scan(&id)
	if err != nil {
		log.Printf("Error inserting product: %v", err)
		return nil, err
	}
	return &pb.CreateProductResponse{Id: fmt.Sprintf("%d", id)}, nil
}

func (s *Server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	var id int
	var name string
	var price float32
	err := s.db.QueryRow("SELECT id, name, price FROM products WHERE id = $1", req.Id).Scan(&id, &name, &price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		log.Printf("Error fetching product: %v", err)
		return nil, err
	}
	return &pb.GetProductResponse{
		Id:    fmt.Sprintf("%d", id),
		Name:  name,
		Price: price,
	}, nil
}
