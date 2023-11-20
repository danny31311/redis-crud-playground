package db

import (
	"context"
	"encoding/json"
	"redis-crud-playground/internals/app/models"

	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type OrdersStorage struct {
	client *redis.Client
}

func NewOrdersStorage(client *redis.Client) *OrdersStorage {
	storage := new(OrdersStorage)
	storage.client = client
	return storage
}

func (storage *OrdersStorage) Create(ctx context.Context, order models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	key := convertToOrderIDKey(order.OrderID)
	tx := storage.client.TxPipeline()

	res := tx.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		tx.Discard()
		logrus.Errorln(err)
		return err
	}

	if err := tx.SAdd(ctx, "orders", key).Err(); err != nil {
		tx.Discard()
		logrus.Errorln(err)
		return err
	}
	if _, err := tx.Exec(ctx); err != nil {
		logrus.Errorln(err)
		return err
	}

	return nil

}
func (storage *OrdersStorage) Update(ctx context.Context, order models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	key := convertToOrderIDKey(order.OrderID)
	res := storage.client.SetXX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil

}

func (storage *OrdersStorage) FindById(ctx context.Context, id uint64) (models.Order, error) {
	res, err := storage.client.Get(ctx, convertToOrderIDKey(id)).Result()
	if err == redis.Nil {
		logrus.Error("Order does not exist ")
		return models.Order{}, err
	} else if err != nil {
		return models.Order{}, err
	}
	var result models.Order
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		logrus.Errorln(err)
	}
	return result, nil

}

func (storage *OrdersStorage) Delete(ctx context.Context, id uint64) error {
	tx := storage.client.TxPipeline()
	key := convertToOrderIDKey(id)
	err := tx.Del(ctx, key).Err()
	if err == redis.Nil {
		tx.Discard()
		logrus.Error("Order does not exist ")
		return err
	} else if err != nil {
		tx.Discard()
		logrus.Errorln(err)
		return err
	}
	if err := tx.SRem(ctx, "orders", key).Err(); err != nil {
		tx.Discard()
		logrus.Errorln(err)
		return err
	}
	if _, err := tx.Exec(ctx); err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}

func (storage *OrdersStorage) List(ctx context.Context, page models.ListPage) (models.FindResult, error) {
	res := storage.client.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		logrus.Errorln(err)
		return models.FindResult{}, err
	}

	if len(keys) == 0 {
		return models.FindResult{Orders: []models.Order{}}, nil
	}

	xs, err := storage.client.MGet(ctx, keys...).Result()
	if err != nil {
		logrus.Errorln(err)
		return models.FindResult{}, err
	}
	orders := make([]models.Order, len(xs))
	for i, x := range xs {
		x := x.(string)
		var order models.Order

		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			logrus.Errorln(err)
			return models.FindResult{}, err
		}
		orders[i] = order
	}
	return models.FindResult{Orders: orders, Cursor: cursor}, nil

}

func convertToOrderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}
