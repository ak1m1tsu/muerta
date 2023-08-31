package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/romankravchuk/muerta/internal/v2/config"
	"github.com/romankravchuk/muerta/internal/v2/lib/logger"
	"github.com/romankravchuk/muerta/internal/v2/services/auth"
	"github.com/romankravchuk/muerta/internal/v2/services/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load()
	failedOnError(err, "failed to load config")

	log := logger.New(cfg.Env)

	log.Debug("config loaded", slog.Any("env", cfg.Env))

	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)

	log.Debug("postgres url created", slog.String("url", postgresURL))

	redisURL := fmt.Sprintf(
		"redis://%s:%s@%s:%s/%s?dial_timeout=%s&read_timeout=%s&write_timeout=%s",
		cfg.Redis.User,
		cfg.Redis.Password,
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Database,
		cfg.Redis.DialTimeout,
		cfg.Redis.ReadTimeout,
		cfg.Redis.WriteTimeout,
	)

	log.Debug("redis url created", slog.String("url", redisURL))

	service, err := auth.New(
		auth.WithAccessCredentials(cfg.Access.PrivateKey, cfg.Access.PublicKey, cfg.Access.Expiration),
		auth.WithRefreshCredentials(cfg.Refresh.PrivateKey, cfg.Refresh.PublicKey, cfg.Refresh.Expiration),
		auth.WithSessionsMemoStorage(),
		auth.WithUsersMemoStorage(),
		auth.WithLogger(log),
	)
	failedOnError(err, "failed to create auth service")

	lis, err := net.Listen("tcp", ":9429")
	failedOnError(err, "failed to create listener")

	log.Info("starting service", slog.String("address", lis.Addr().String()))

	gsrv := grpc.NewServer()

	proto.RegisterAuthServiceServer(gsrv, service)
	reflection.Register(gsrv)

	go func() {
		failedOnError(gsrv.Serve(lis), "failed to start service")
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Info("stopping service")

	gsrv.GracefulStop()
	os.Exit(0)
}

func failedOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, slog.String("error", err.Error()))
		os.Exit(1)
	}
}
