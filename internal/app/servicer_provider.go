package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/merynayr/auth/internal/closer"
	"github.com/merynayr/auth/internal/config"
	"github.com/merynayr/auth/internal/config/env"
	"github.com/merynayr/auth/internal/repository"
	"github.com/merynayr/auth/internal/service"

	userAPI "github.com/merynayr/auth/internal/api/user"
	userRepository "github.com/merynayr/auth/internal/repository/user"
	userService "github.com/merynayr/auth/internal/service/user"
)

// Структура приложения со всеми зависимости
type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgPool *pgxpool.Pool

	userRepository repository.UserRepository
	userService    service.UserService
	userAPI        *userAPI.API
}

// NewServiceProvider возвращает новый объект API слоя
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get gprc config: %s", err.Error())
		}
		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		if err := pool.Ping(ctx); err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.pgPool = pool
	}

	return s.pgPool
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.PgPool(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserAPI(ctx context.Context) *userAPI.API {
	if s.userAPI == nil {
		s.userAPI = userAPI.NewAPI(s.UserService(ctx))
	}

	return s.userAPI
}
