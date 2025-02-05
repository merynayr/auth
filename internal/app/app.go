package app

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/merynayr/auth/internal/closer"
	"github.com/merynayr/auth/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	descUser "github.com/merynayr/auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// App структура приложения
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp возвращает объект приложения
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает приложение
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

// initDeps инициализирует все зависимости
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	descUser.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserAPI(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())
	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
