package services

import (
	"context"
	"errors"
	"redis-crud-playground/internals/app/db"
	"redis-crud-playground/internals/app/models"
)

type OrdersService struct {
	storage *db.OrdersStorage
}

func NewOrdersService(storage *db.OrdersStorage) *OrdersService {
	service := new(OrdersService)
	service.storage = storage
	return service
}

func (service *OrdersService) CreateOrder(ctx context.Context, order models.Order) error {
	return service.storage.Create(ctx, order)
}

func (service *OrdersService) UpdateOrder(ctx context.Context, order models.Order) error {
	return service.storage.Update(ctx, order)
}
func (service *OrdersService) DeleteOrder(ctx context.Context, id uint64) error {
	return service.storage.Delete(ctx, id)
}

func (service *OrdersService) FindById(ctx context.Context, id uint64) (models.Order, error) {
	order, err := service.storage.FindById(ctx, id)
	if err != nil {
		return order, errors.New("order not found")
	}
	return order, nil
}

func (service *OrdersService) List(ctx context.Context, page models.ListPage) (models.FindResult, error) {
	return service.storage.List(ctx, page)
}
