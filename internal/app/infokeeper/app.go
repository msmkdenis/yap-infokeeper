package infokeeper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/msmkdenis/yap-infokeeper/internal/config"
	grpcCredentialHandlers "github.com/msmkdenis/yap-infokeeper/internal/credential/api/grpchandlers"
	pbCredential "github.com/msmkdenis/yap-infokeeper/internal/credential/api/grpchandlers/proto"
	credentialRepository "github.com/msmkdenis/yap-infokeeper/internal/credential/repository"
	credentialService "github.com/msmkdenis/yap-infokeeper/internal/credential/service"
	grpcCreditCardHandlers "github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers"
	pbCreditCard "github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers/proto"
	creditCardRepository "github.com/msmkdenis/yap-infokeeper/internal/credit_card/repository"
	creditCardService "github.com/msmkdenis/yap-infokeeper/internal/credit_card/service"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
	grpcTextDataHandlers "github.com/msmkdenis/yap-infokeeper/internal/text_data/api/grpchandlers"
	pbTextData "github.com/msmkdenis/yap-infokeeper/internal/text_data/api/grpchandlers/proto"
	textDataRepository "github.com/msmkdenis/yap-infokeeper/internal/text_data/repository"
	textDataService "github.com/msmkdenis/yap-infokeeper/internal/text_data/service"
	grpcUserHandlers "github.com/msmkdenis/yap-infokeeper/internal/user/api/grpchandlers"
	pbUser "github.com/msmkdenis/yap-infokeeper/internal/user/api/grpchandlers/proto"
	userRepository "github.com/msmkdenis/yap-infokeeper/internal/user/repository"
	userService "github.com/msmkdenis/yap-infokeeper/internal/user/service"
	"github.com/msmkdenis/yap-infokeeper/pkg/jwtgen"
)

func Run(quitSignal chan os.Signal) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)
	cfg := config.New()

	jwtManager := jwtgen.NewJWTManager(cfg.TokenName, cfg.TokenSecret, cfg.TokenExpHours)
	jwtAuth := interceptors.NewJWTAuth(jwtManager)

	postgresPool := initPostgresPool(cfg)

	userRepo := userRepository.NewPostgresUserRepository(postgresPool)
	userServ := userService.NewUserService(userRepo)

	creditCardRepo := creditCardRepository.NewPostgresCreditCardRepository(postgresPool)
	creditCardServ := creditCardService.NewCreditCardService(creditCardRepo)

	credentialRepo := credentialRepository.NewPostgresCredentialsRepository(postgresPool)
	credentialServ := credentialService.NewCredentialService(credentialRepo)

	textDataRepo := textDataRepository.NewPostgresTextDataRepository(postgresPool)
	textDataServ := textDataService.NewTextDataService(textDataRepo)

	listener, err := net.Listen("tcp", cfg.GRPCServer)
	if err != nil {
		slog.Error("Unable to create listener", slog.String("error", err.Error()))
		os.Exit(1)
	}
	serverGrpc := grpc.NewServer(
		grpc.ChainUnaryInterceptor(jwtAuth.GRPCJWTAuth),
	)

	pbUser.RegisterUserServiceServer(serverGrpc, grpcUserHandlers.NewUserRegister(userServ, jwtManager))
	pbCreditCard.RegisterCreditCardServiceServer(serverGrpc, grpcCreditCardHandlers.NewCreditCard(creditCardServ))
	pbCredential.RegisterCredentialServiceServer(serverGrpc, grpcCredentialHandlers.NewCredential(credentialServ))
	pbTextData.RegisterTextDataServiceServer(serverGrpc, grpcTextDataHandlers.NewTextData(textDataServ))

	reflection.Register(serverGrpc)

	grpcServerCtx, grpcServerStopCtx := context.WithCancel(context.Background())

	quit := make(chan struct{})
	go func() {
		<-quitSignal
		close(quit)
	}()

	go func() {
		logger.Info(fmt.Sprintf("gRPC server starting on port %s", cfg.GRPCServer))
		if err := serverGrpc.Serve(listener); err != nil {
			slog.Error("Unable to start gRPC server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	go func() {
		<-quit

		// Shutdown signal with grace period of 10 seconds
		shutdownCtx, cancel := context.WithTimeout(grpcServerCtx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				logger.Error("graceful gRPC shutdown timed out.. forcing exit.")
				serverGrpc.Stop()
			}
		}()

		// Trigger graceful shutdown
		logger.Info("Shutdown signal received, gracefully stopping gRPC server")
		serverGrpc.GracefulStop()
		grpcServerStopCtx()
	}()

	<-grpcServerCtx.Done()
}

func initPostgresPool(cfg *config.Config) *db.PostgresPool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	postgresPool, err := db.NewPostgresPool(ctx, cfg.DatabaseURI)
	if err != nil {
		slog.Error("Unable to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	migrations, err := db.NewMigrations(cfg.DatabaseURI)
	if err != nil {
		slog.Error("Unable to create migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = migrations.MigrateUp()
	if err != nil {
		slog.Error("Unable to up migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer cancel()
	slog.Info("Connected to database", slog.String("DSN", cfg.DatabaseURI))
	return postgresPool
}
