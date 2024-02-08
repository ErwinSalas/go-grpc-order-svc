package order

import (
	"context"
	"net/http"

	"github.com/ErwinSalas/go-grpc-order-svc/pkg/client"
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/models"
	orderpb "github.com/ErwinSalas/go-grpc-order-svc/proto"
)

type OrderService interface {
	CreateOrder(context.Context, *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error)
}

type orderService struct {
	repository    OrderRepository
	productClient client.ProductServiceClient
}

func NewOrderService(repository OrderRepository, productClient client.ProductServiceClient) OrderService {
	return &orderService{
		repository,
		productClient,
	}
}

// CreateOrder handles the logic for creating a new order.
func (s *orderService) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	product, err := s.productClient.FindOne(req.ProductId)

	if err != nil {
		return &orderpb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	} else if product.Status >= http.StatusNotFound {
		return &orderpb.CreateOrderResponse{Status: product.Status, Error: product.Error}, nil
	} else if product.Data.Stock < req.Quantity {
		return &orderpb.CreateOrderResponse{Status: http.StatusConflict, Error: "Stock too low"}, nil
	}

	order := models.Order{
		Price:     product.Data.Price,
		ProductId: product.Data.Id,
		UserId:    req.UserId,
	}

	err = s.repository.CreateOrder(&order)

	if err != nil {
		return &orderpb.CreateOrderResponse{Status: http.StatusInternalServerError, Error: "Failed to create order"}, nil
	}

	res, err := s.productClient.DecreaseStock(req.ProductId, order.Id)

	if err != nil {
		return &orderpb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	} else if res.Status == http.StatusConflict {
		s.repository.DeleteOrder(order.Id)

		return &orderpb.CreateOrderResponse{Status: http.StatusConflict, Error: res.Error}, nil
	}

	return &orderpb.CreateOrderResponse{
		Status: http.StatusCreated,
		Id:     order.Id,
	}, nil
}
