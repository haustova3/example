package app

import (
	"context"
	"example/comments/internal/app/config"
	"example/comments/internal/app/middlewares"
	"example/comments/internal/external/notification"
	"example/comments/internal/external/products"
	"example/comments/internal/external/users"
	"example/comments/internal/logger"
	"example/comments/internal/repository"
	"example/comments/internal/trace"
	"example/comments/internal/usecases"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	desc "example/comments/pkg/api/comments/v1"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config        *config.Config
	rep           *repository.Repository
	grpcServer    *grpc.Server
	gwServer      *http.Server
	metricsServer *http.Server
}

func NewApp(ctx context.Context, configPath string) (*App, error) {
	configImpl, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{
		config: configImpl,
	}
	appCtx, cancel := context.WithCancel(ctx)
	app.ConnectDatabase(appCtx, configImpl.NotificationConf.MaxCount)
	notification.StartNotificationService(appCtx, app.rep,
		[]string{app.config.KafkaConf.Brokers},
		app.config.KafkaConf.OrderTopic,
		app.config.NotificationConf.MaxCount,
		app.config.NotificationConf.Timer)
	app.SignalHandler(ctx, cancel)
	trace.CreateTracerProvider(appCtx, configImpl)
	return app, nil
}

func (app *App) ConnectDatabase(appCtx context.Context, ntfMaxCount int) {
	address := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		app.config.DBConf.User, app.config.DBConf.Password, app.config.DBConf.Host, app.config.DBConf.Port, app.config.DBConf.DBName)
	conf, err := pgxpool.ParseConfig(address)
	if err != nil {
		logger.Errorw(appCtx, "unable to parse master repository config", "err", err)
		panic(err)
	}
	masterPool, err := pgxpool.NewWithConfig(appCtx, conf)
	if err != nil {
		logger.Errorw(appCtx, "unable to create pgx pool master", "err", err)
		panic(err)
	}

	app.rep = repository.NewRepository(masterPool, ntfMaxCount)
}

func (app *App) ListenAndServe(ctx context.Context) error {
	address := fmt.Sprintf(":%s", app.config.ServiceConf.GRPCPort)
	logger.Infow(ctx, "Starting grpc server", "address", address)
	list, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen and serve app failed: %w", err)
	}

	app.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.Panic,
			mw.Logger,
			mw.Trace,
			mw.Validate,
		),
	)

	reflection.Register(app.grpcServer)

	productsAddress := fmt.Sprintf("%s:%s", app.config.ProductsConf.Host, app.config.ProductsConf.Port)
	productsService, err := products.NewProductsService(ctx, productsAddress)
	if err != nil {
		return err
	}

	usersAddress := fmt.Sprintf("%s:%s", app.config.UsersConf.Host, app.config.UsersConf.Port)
	usersService, err := users.NewUsersService(ctx, usersAddress)
	if err != nil {
		return err
	}

	createCommentService := usecases.NewCreateCommentService(app.rep, productsService, usersService)
	getCommentsService := usecases.NewGetCommentsService(app.rep)
	commentsController := NewCommentsController(createCommentService, getCommentsService)
	desc.RegisterCommentsServer(app.grpcServer, commentsController)

	logger.Infow(ctx, "server listening", "address", list.Addr())
	go func() {
		if err = app.grpcServer.Serve(list); err != nil {
			panic(err)
		}
	}()

	go func() {
		address := fmt.Sprintf("%s:%s", app.config.ServiceConf.Host, app.config.ServiceConf.MetricPort)
		logger.Infow(context.Background(), "Starting metric loms service", "address", address)
		l, _ := net.Listen("tcp", address)
		mx := http.NewServeMux()
		mx.Handle("GET /metrics", promhttp.Handler())
		app.metricsServer = &http.Server{}
		app.metricsServer.Handler = mx
		if err = app.metricsServer.Serve(l); err != nil {
			panic(err)
		}
	}()

	err = app.CreateHTTPGateway(ctx)
	if err != nil {
		logger.Errorw(ctx, "Can not start HTTP gateway", "err", err)
		return err
	}

	return nil
}

func (app *App) CreateHTTPGateway(ctx context.Context) error {
	address := fmt.Sprintf("%s:%s", app.config.ServiceConf.Host, app.config.ServiceConf.HTTPPort)
	logger.Infow(ctx, "Starting gateway", "address", address)
	targetAddress := fmt.Sprintf(":%s", app.config.ServiceConf.GRPCPort)
	conn, err := grpc.NewClient(targetAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorw(ctx, "Failed to deal", "err", err)
		panic(err)
	}

	gwmux := runtime.NewServeMux()

	if err = desc.RegisterCommentsHandler(context.Background(), gwmux, conn); err != nil {
		logger.Errorw(ctx, "Failed to register gateway", "err", err)
		panic(err)
	}

	app.gwServer = &http.Server{
		Addr:              address,
		Handler:           gwmux,
		ReadHeaderTimeout: time.Duration(1) * time.Second,
		ReadTimeout:       time.Duration(3) * time.Second,
	}

	logger.Infow(ctx, "Serving gRPC-Gateway", "address", app.gwServer.Addr)
	if err = app.gwServer.ListenAndServe(); err != nil {
		panic(err)
	}

	return nil
}

func (app *App) SignalHandler(ctx context.Context, appCancelContext func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func(appCancelContext func()) {
		sig := <-sigChan
		sdCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		logger.Infow(sdCtx, "Received signal", "signal", sig)
		err := app.gwServer.Shutdown(sdCtx)
		if err != nil {
			logger.Infow(sdCtx, "can not shutdown server", "err", err.Error())
		}
		app.grpcServer.GracefulStop()
		appCancelContext()
	}(appCancelContext)
}
