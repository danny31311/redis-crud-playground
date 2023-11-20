package handlers

import (
	"encoding/json"
	"errors"
	"redis-crud-playground/internals/app/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"redis-crud-playground/internals/app/services"
	"strconv"
	"time"
)

type OrdersHandler struct {
	service *services.OrdersService
}

func NewOrdersHandler(service *services.OrdersService) *OrdersHandler {
	handler := new(OrdersHandler)
	handler.service = service
	return handler
}

func (handler *OrdersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID uuid.UUID         `json:"customer_id"`
		LineItems  []models.LineItem `json:"line_items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WrapError(w, err)
		return
	}
	now := time.Now().UTC()

	order := models.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreateAt:   &now,
	}
	err := handler.service.CreateOrder(r.Context(), order)
	if err != nil {
		WrapError(w, err)
		return
	}
	res, err := json.Marshal(order)
	if err != nil {
		WrapError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (handler *OrdersHandler) List(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	if cursor == "" {
		cursor = "0"
	}

	curs, err := strconv.ParseUint(cursor, 10, 64)
	if err != nil {
		WrapError(w, err)
		return
	}
	const size = 30
	res, err := handler.service.List(r.Context(), models.ListPage{
		Offset: curs,
		Size:   size,
	})
	if err != nil {
		WrapError(w, err)
		return
	}
	var response struct {
		Items []models.Order `json:"items"`
		Next  uint64         `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		WrapError(w, err)
		return
	}
	w.Write(data)

}

func (handler *OrdersHandler) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		WrapError(w, errors.New("missing id"))
		return
	}
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		WrapError(w, err)
		return
	}
	order, err := handler.service.FindById(r.Context(), id)
	if err != nil {
		WrapError(w, err)
		return
	}
	if err := json.NewEncoder(w).Encode(order); err != nil {
		WrapError(w, err)
		return

	}

}

func (handler *OrdersHandler) UpdateById(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WrapError(w, err)
		return
	}
	vars := mux.Vars(r)
	if vars["id"] == "" {
		WrapError(w, errors.New("missing id"))
		return
	}
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		WrapError(w, err)
		return
	}
	order, err := handler.service.FindById(r.Context(), id)
	if err != nil {
		WrapError(w, err)
		return
	}

	now := time.Now().UTC()

	switch body.Status {
	case "shipped":
		if order.ShippedAt != nil {
			WrapError(w, err)
			return
		}
		order.ShippedAt = &now
	case "completed":
		if order.CompletedAt != nil || order.ShippedAt == nil {
			WrapError(w, err)
			return
		}
		order.CompletedAt = &now
	default:
		WrapError(w, err)
		return
	}
	err = handler.service.UpdateOrder(r.Context(), order)
	if err != nil {
		WrapError(w, err)
		return
	}
	if err := json.NewEncoder(w).Encode(order); err != nil {
		WrapError(w, err)
		return
	}
}

func (handler *OrdersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		WrapError(w, errors.New("missing id"))
		return
	}
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		WrapError(w, err)
		return
	}
	err = handler.service.DeleteOrder(r.Context(), id)
	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"result": "OK",
		"data":   "",
	}
	WrapOk(w, m)
}
