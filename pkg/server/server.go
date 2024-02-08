package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ErwinSalas/go-grpc-order-svc/pkg/order"
	orderpb "github.com/ErwinSalas/go-grpc-order-svc/proto"
)

type orderServer struct {
	orderpb.UnimplementedOrderServiceServer
	order.OrderService
}

func NewOrderServer(orderService order.OrderService) orderpb.OrderServiceServer {
	return &orderServer{
		orderpb.UnimplementedOrderServiceServer{},
		orderService,
	}
}

func (s *orderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	fmt.Println("RPC order-service/CreateOrder")

	resp, err := s.OrderService.CreateOrder(ctx, req)

	if err != nil {
		return &orderpb.CreateOrderResponse{
			Status: resp.Status,
		}, err
	}

	return &orderpb.CreateOrderResponse{
		Status: http.StatusCreated,
		Id:     resp.Id,
	}, nil
}
