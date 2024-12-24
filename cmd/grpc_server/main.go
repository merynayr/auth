package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	desc "github.com/merynayr/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/merynayr/auth/internal/config"
	"github.com/merynayr/auth/internal/config/env"
)

type server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	op := "GetUser"
	log.Printf("[%s] request data | id: %v", op, req.Id)
	builderGet := sq.Select("id", "username", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderGet.ToSql()
	if err != nil {
		log.Printf("%s: failed to create builder: %v", op, err)
		return nil, err
	}
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("%s: failed to select user: %v", op, err)
		return nil, err
	}

	var userID int64
	var username, email string
	var role int32
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&userID, &username, &email, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("%s: failed to scan user: %v", op, err)
			return nil, err
		}
	}

	log.Printf("%s: selected user %d", op, req.Id)
	return &desc.GetResponse{
		Auth: &desc.Auth{
			Id:        userID,
			Name:      username,
			Email:     email,
			Role:      desc.Role(role),
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt.Time),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	op := "CreateUser"
	log.Printf("%s %v", op, req.GetInfo())

	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("username", "email", "password", "role", "created_at").
		Values(req.Info.Name, req.Info.Email, req.Info.Password, req.Info.GetRole(), time.Now()).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("%s: failed to create builder: %v", op, err)
		return nil, err
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("%s: failed to insert user: %v", op, err)
		return nil, err
	}

	log.Printf("%s: inserted user with id: %d", op, userID)
	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
