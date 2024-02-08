package order

import (
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/database"
	"github.com/ErwinSalas/go-grpc-order-svc/pkg/models"
)

// OrderRepository is the interface for manipulating orders in the database.
type OrderRepository interface {
	CreateOrder(order *models.Order) error
	DeleteOrder(orderID int64) error
}

// GormOrderRepository is the Gorm implementation of the OrderRepository interface.
type gormOrderRepository struct {
	datastore database.Handler
}

func NewOrderRepository(datastore database.Handler) OrderRepository {
	return &gormOrderRepository{
		datastore,
	}
}
func (r *gormOrderRepository) CreateOrder(order *models.Order) error {
	return r.datastore.DB.Create(order).Error
}

func (r *gormOrderRepository) DeleteOrder(orderID int64) error {
	return r.datastore.DB.Delete(&models.Order{}, orderID).Error
}
