package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	authInterceptor "github.com/ffauzann/loan-service/pkg/common/interceptor/grpc/unary/authentication"
	ctxTagsInterceptor "github.com/ffauzann/loan-service/pkg/common/interceptor/grpc/unary/ctxtags"
	logInterceptor "github.com/ffauzann/loan-service/pkg/common/interceptor/grpc/unary/logging"
	recoveryInterceptor "github.com/ffauzann/loan-service/pkg/common/interceptor/grpc/unary/recovery"

	grpcCtxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/ffauzann/loan-service/client"
	"github.com/ffauzann/loan-service/internal/constant"
	deliveryGRPC "github.com/ffauzann/loan-service/internal/delivery/grpc"
	deliveryHTTP "github.com/ffauzann/loan-service/internal/delivery/http"
	deliveryKafka "github.com/ffauzann/loan-service/internal/delivery/kafka"
	"github.com/ffauzann/loan-service/internal/repository"
	"github.com/ffauzann/loan-service/internal/service"
	"github.com/ffauzann/loan-service/internal/util"
	"github.com/ffauzann/loan-service/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

type Server struct {
	GRPC   GRPC
	HTTP   HTTP
	Logger Logger
}

type GRPC struct {
	Address string
	Port    uint32
	Server  *grpc.Server
}

type HTTP struct {
	Address string
	Port    uint32
	Timeout string
	Server  *http.Server
}

// StartServer initializes and starts the gRPC and HTTP servers, sets up the repositories, services, and handles graceful shutdown.
// It also starts a Kafka consumer to process messages from the messaging system.
// The function uses goroutines to run the servers concurrently and waits for an interrupt signal to gracefully shut down the servers.
func (c *Config) StartServer() {
	var wg sync.WaitGroup
	wg.Add(3) //nolint

	// Init repo
	dbRepo := repository.NewDB(c.Database.SQL.DB, c.App, c.Server.Logger.Zap)
	redisRepo := repository.NewRedis(c.Cache.Redis.Client, c.App, c.Server.Logger.Zap)
	messagingRepo := repository.NewMessaging(c.Messaging.Kafka.Producer, c.App, c.Server.Logger.Zap)
	notifRepo := repository.NewNotification(c.SMTP.MailHog.Client, c.App, c.Server.Logger.Zap)

	// Init service
	svc := service.New(dbRepo, redisRepo, messagingRepo, notifRepo, c.App, c.Server.Logger.Zap)

	go func() {
		defer wg.Done()
		c.startGRPCServer(svc)
	}()

	go func() {
		defer wg.Done()
		c.startHTTPProxyServer(svc)
	}()

	go func() {
		defer wg.Done()
		c.startConsumer(svc)
	}()

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	fmt.Println("Gracefully shutting down...")

	c.Server.GRPC.Server.GracefulStop()
	fmt.Println("gRPC server has been shutdown.")
	c.Server.HTTP.Server.Shutdown(context.Background())
	fmt.Println("HTTP proxy server has been shutdown.")
	c.Messaging.Kafka.Producer.Close()
	fmt.Println("Kafka producer has been shutdown.")
	c.Messaging.Kafka.Consumer.Close()
	fmt.Println("Kafka consumer has been shutdown.")

}

// startGRPCServer starts the gRPC server and registers the service handlers.
// It listens on the specified address and port, and sets up interceptors for logging, authentication, and recovery.
func (c *Config) startGRPCServer(svc service.Service) {
	addr := fmt.Sprintf("%s:%d", c.Server.GRPC.Address, c.Server.GRPC.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	interceptors := []grpc.UnaryServerInterceptor{
		ctxTagsInterceptor.UnaryServerInterceptor(),
		recoveryInterceptor.UnaryServerInterceptor(c.Server.Logger.Zap),
		authInterceptor.UnaryServerInterceptor(
			&authInterceptor.Config{
				Iss:     c.App.Jwt.AccessToken.Iss,
				Alg:     authInterceptor.AlgRS256,
				JwksURL: fmt.Sprintf("http://%s:%d%s", c.Server.HTTP.Address, c.Server.HTTP.Port, client.JwksPath),
			},
			authInterceptor.WithCustomMetadataKey("Authorization"),
			authInterceptor.WithCustomClaims(&client.Claims{}),
			authInterceptor.WithExcludedMethods(c.App.Auth.ExcludedMethods...),
		),
		logInterceptor.UnaryServerInterceptor(
			c.Server.Logger.Zap,
			logInterceptor.WithErrorParser(constant.MapGRPCErrCodes),
			logInterceptor.WithCtxTag(true),
		),
		grpcCtxtags.UnaryServerInterceptor(),
	}
	opts := grpc.ChainUnaryInterceptor(interceptors...)
	c.Server.GRPC.Server = grpc.NewServer(opts)

	deliveryGRPC.New(c.Server.GRPC.Server, svc, c.Server.Logger.Zap)
	fmt.Printf("gRPC server started on %s\n", addr)

	if err := c.Server.GRPC.Server.Serve(lis); err != nil {
		log.Fatal(err)
		return
	}
}

// startHTTPProxyServer starts an HTTP proxy server that forwards requests to the gRPC server.
// It uses the grpc-gateway to handle HTTP requests and convert them to gRPC calls.
func (c *Config) startHTTPProxyServer(svc service.Service) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcMux := runtime.NewServeMux(
		runtime.WithMetadata(
			func(ctx context.Context, r *http.Request) metadata.MD {
				return metadata.Pairs("X-Forwarded-Method", r.Method)
			},
		),
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:     true,
					EmitDefaultValues: true,
				},
			},
		),
	)

	timeout, _ := time.ParseDuration(c.Server.HTTP.Timeout)
	addr := fmt.Sprintf("%s:%d", c.Server.HTTP.Address, c.Server.HTTP.Port)
	c.Server.HTTP.Server = &http.Server{
		Addr:              addr,
		Handler:           grpcMux,
		ReadHeaderTimeout: timeout,
		ReadTimeout:       timeout,
		WriteTimeout:      timeout,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := gen.RegisterAuthServiceHandlerFromEndpoint(ctx, grpcMux, fmt.Sprintf("%s:%d", c.Server.GRPC.Address, c.Server.GRPC.Port), opts); err != nil {
		util.Log().Error(err.Error())
		return
	}

	if err := gen.RegisterUserServiceHandlerFromEndpoint(ctx, grpcMux, fmt.Sprintf("%s:%d", c.Server.GRPC.Address, c.Server.GRPC.Port), opts); err != nil {
		util.Log().Error(err.Error())
		return
	}

	if err := gen.RegisterLoanServiceHandlerFromEndpoint(ctx, grpcMux, fmt.Sprintf("%s:%d", c.Server.GRPC.Address, c.Server.GRPC.Port), opts); err != nil {
		util.Log().Error(err.Error())
		return
	}

	deliveryHTTP.New(grpcMux, svc)

	fmt.Printf("HTTP proxy server started on %s\n", addr)

	if err := c.Server.HTTP.Server.ListenAndServe(); err != nil {
		util.Log().Error(err.Error())
		return
	}
}

// startConsumer starts a Kafka consumer that listens for messages and processes them accordingly.
func (c *Config) startConsumer(svc service.Service) {
	consumer := c.Messaging.Kafka.Consumer
	deliveryKafka.New(consumer, svc)

	// for {
	// 	msg, err := reader.ReadMessage(context.Background())
	// 	if err != nil {
	// 		c.Server.Logger.Zap.Error("failed to read kafka msg", zap.Error(err))
	// 		continue
	// 	}

	// 	switch msg.Topic {
	// 	case constant.TopicFullyInvested:
	// 		go func() {
	// 			if err := svc.Notification.SendEmail(msg.Value); err != nil {
	// 				c.Server.Logger.Zap.Error("email failed", zap.Error(err))
	// 			}
	// 		}()
	// 	case "loan-service.loan-approved":
	// 		go func() {
	// 			if err := svc.Loan.ProcessApproved(msg.Value); err != nil {
	// 				c.Server.Logger.Zap.Error("loan process failed", zap.Error(err))
	// 			}
	// 		}()
	// 	default:
	// 		c.Server.Logger.Zap.Warn("unhandled kafka topic", zap.String("topic", msg.Topic))
	// 	}
	// }
}
