package client

import (
	"context"
	"fmt"

	productpb "github.com/ErwinSalas/go-grpc-product-svc/proto"
	"google.golang.org/grpc"
)

type ProductServiceClient struct {
	Client productpb.ProductServiceClient
}

func NewProductServiceClient(url string) ProductServiceClient {
	cc, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := ProductServiceClient{
		Client: productpb.NewProductServiceClient(cc),
	}

	return c
}

func (c *ProductServiceClient) FindOne(productId int64) (*productpb.FindOneResponse, error) {
	req := &productpb.FindOneRequest{
		Id: productId,
	}

	return c.Client.FindOne(context.Background(), req)
}

func (c *ProductServiceClient) DecreaseStock(productId int64, orderId int64) (*productpb.DecreaseStockResponse, error) {
	req := &productpb.DecreaseStockRequest{
		Id:      productId,
		OrderId: orderId,
	}

	return c.Client.DecreaseStock(context.Background(), req)
}
