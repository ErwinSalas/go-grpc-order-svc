package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ErwinSalas/go-grpc-order-svc/pkg/client"
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/config"
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/database"
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/order"
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/server"
	orderpb "github.com/ErwinSalas/go-grpc-order-svc/proto"
	"google.golang.org/grpc"
)

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

	grpcServer := grpc.NewServer()

	orderpb.RegisterOrderServiceServer(grpcServer, server.NewOrderServer(orderService))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

}
