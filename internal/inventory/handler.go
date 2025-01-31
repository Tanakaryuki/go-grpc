package inventory

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	pb "github.com/Tanakaryuki/go-grpc/api/inventory"
	_ "github.com/lib/pq"
)

type Server struct {
	pb.UnimplementedInventoryServiceServer
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
    CREATE TABLE IF NOT EXISTS inventory (
        id SERIAL PRIMARY KEY,
        product_id INTEGER UNIQUE,
        quantity INTEGER
    );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &Server{db: db}, nil
}

func (s *Server) AddInventory(ctx context.Context, req *pb.AddInventoryRequest) (*pb.AddInventoryResponse, error) {
	var id int
	// Upsert操作
	err := s.db.QueryRow(`
        INSERT INTO inventory (product_id, quantity) VALUES ($1, $2)
        ON CONFLICT (product_id) DO UPDATE SET quantity = inventory.quantity + EXCLUDED.quantity
        RETURNING id
    `, req.ProductId, req.Quantity).Scan(&id)
	if err != nil {
		log.Printf("Error adding inventory: %v", err)
		return nil, err
	}
	return &pb.AddInventoryResponse{Id: fmt.Sprintf("%d", id)}, nil
}

func (s *Server) GetInventory(ctx context.Context, req *pb.GetInventoryRequest) (*pb.GetInventoryResponse, error) {
	var productId, quantity int
	err := s.db.QueryRow("SELECT product_id, quantity FROM inventory WHERE product_id = $1", req.ProductId).
		Scan(&productId, &quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("inventory not found")
		}
		log.Printf("Error fetching inventory: %v", err)
		return nil, err
	}
	return &pb.GetInventoryResponse{
		ProductId: fmt.Sprintf("%d", productId),
		Quantity:  int32(quantity),
	}, nil
}
