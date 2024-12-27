package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	desc "github.com/merynayr/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/merynayr/auth/internal/config"
	"github.com/merynayr/auth/internal/config/env"
	"github.com/merynayr/auth/internal/converter"
	"github.com/merynayr/auth/internal/repository"
	"github.com/merynayr/auth/internal/repository/user"
)

type server struct {
	desc.UnimplementedUserV1Server
	userRepository repository.UserRepository
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

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping to database: %v", err)
	}

	userRepo := user.NewRepository(pool)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{userRepository: userRepo})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Create(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	id, err := s.userRepository.CreateUser(ctx, converter.ToUserFromDescUser(req))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted note with id: %d", id)

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	userObj, err := s.userRepository.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToDescUserFromService(userObj), nil

}
