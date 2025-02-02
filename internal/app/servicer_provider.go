package app

import (
	"context"
	"log"

	"github.com/merynayr/auth/internal/client/db"
	"github.com/merynayr/auth/internal/client/db/pg"
	"github.com/merynayr/auth/internal/client/db/transaction"
	"github.com/merynayr/auth/internal/closer"
	"github.com/merynayr/auth/internal/config"
	"github.com/merynayr/auth/internal/config/env"
	"github.com/merynayr/auth/internal/repository"
	"github.com/merynayr/auth/internal/service"

	accessAPI "github.com/merynayr/auth/internal/api/access"
	authAPI "github.com/merynayr/auth/internal/api/auth"
	userAPI "github.com/merynayr/auth/internal/api/user"

	accessService "github.com/merynayr/auth/internal/service/access"
	authService "github.com/merynayr/auth/internal/service/auth"
	userService "github.com/merynayr/auth/internal/service/user"

	userRepository "github.com/merynayr/auth/internal/repository/user"
)

// Структура приложения со всеми зависимости
type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig
	authConfig    config.AuthConfig
	accessConfig  config.AccessConfig

	dbClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository
	userService    service.UserService
	userAPI        *userAPI.API

	authService service.AuthService
	authAPI     *authAPI.API

	accessService service.AccessService
	accessAPI     *accessAPI.API
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

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

// AuthConfig инициализирует конфиг auth сервиса
func (s *serviceProvider) AuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := env.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get auth config")
		}

		s.authConfig = cfg
	}

	return s.authConfig
}

// AccessConfig инициализирует конфиг access конфига
func (s *serviceProvider) AccessConfig() config.AccessConfig {
	if s.accessConfig == nil {
		cfg, err := env.NewAccessConfig()
		if err != nil {
			log.Fatalf("failed to get access service")
		}

		s.accessConfig = cfg
	}

	return s.accessConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

// AuthService иницилизирует сервисный слой auth
func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.UserRepository(ctx),
			s.AuthConfig(),
		)
	}

	return s.authService
}

// AccessService иницилизирует сервисный слой access
func (s *serviceProvider) AccessService(_ context.Context) service.AccessService {
	if s.accessService == nil {
		uMap, err := s.AccessConfig().UserAccessesMap()
		if err != nil {
			log.Fatalf("failed to get user access map: %v", err)
		}

		s.accessService = accessService.NewService(uMap, s.AuthConfig())
	}

	return s.accessService
}

func (s *serviceProvider) UserAPI(ctx context.Context) *userAPI.API {
	if s.userAPI == nil {
		s.userAPI = userAPI.NewAPI(s.UserService(ctx))
	}

	return s.userAPI
}

// AuthAPI инициализирует api слой auth
func (s *serviceProvider) AuthAPI(ctx context.Context) *authAPI.API {
	if s.authAPI == nil {
		s.authAPI = authAPI.NewAPI(s.AuthService(ctx))
	}

	return s.authAPI
}

// AccessAPI инициализирует api слой access
func (s *serviceProvider) AccessAPI(ctx context.Context) *accessAPI.API {
	if s.accessAPI == nil {
		s.accessAPI = accessAPI.NewAPI(s.AccessService(ctx))
	}

	return s.accessAPI
}
