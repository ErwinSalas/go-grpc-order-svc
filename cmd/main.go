package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ErwinSalas/go-grpc-order-svc/pkg/client"
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/config"
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/database"
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/order"
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/server"
	orderpb "github.com/ErwinSalas/go-grpc-order-svc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Received request: %v", req)
	log.Printf("Method: %s", info.FullMethod)

	resp, err := handler(ctx, req)
	return resp, err
}

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Order Svc on", c.Port)

	datastore := database.Init(c.DBUrl)
	productClient := client.NewProductServiceClient(c.ProductSvcUrl)
	orderRepository := order.NewOrderRepository(datastore)
	orderService := order.NewOrderService(orderRepository, productClient)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))

	orderpb.RegisterOrderServiceServer(grpcServer, server.NewOrderServer(orderService))

	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthcheck)

	// Start health check routine
	go func() {
		for {
			var count int64
			if err := datastore.DB.Table("orders").Count(&count).Error; err != nil {
				log.Println("Database query error:", err)
				healthcheck.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
				return
			} else {
				healthcheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
				time.Sleep(5 * time.Second)
			}
		}
	}()
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

}
