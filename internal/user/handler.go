package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	pb "github.com/Tanakaryuki/go-grpc/api/user"
	_ "github.com/lib/pq"
)

type Server struct {
	pb.UnimplementedUserServiceServer
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
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        email VARCHAR(100) UNIQUE
    );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &Server{db: db}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	var id int
	err := s.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		req.Name, req.Email).Scan(&id)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return nil, err
	}
	return &pb.CreateUserResponse{Id: fmt.Sprintf("%d", id)}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var id int
	var name, email string
	err := s.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", req.Id).Scan(&id, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	return &pb.GetUserResponse{
		Id:    fmt.Sprintf("%d", id),
		Name:  name,
		Email: email,
	}, nil
}
