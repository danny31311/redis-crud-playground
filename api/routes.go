package api

import (
	"github.com/gorilla/mux"
	"redis-crud-playground/internals/app/handlers"
)

func CreateRoutes(orderHandler *handlers.OrdersHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", orderHandler.Create).Methods("POST")
	r.HandleFunc("/", orderHandler.List).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", orderHandler.FindById).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", orderHandler.UpdateById).Methods("PUT")
	r.HandleFunc("/{id:[0-9]+}", orderHandler.Delete).Methods("DELETE")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(handlers.NotFound).GetHandler()

	return r
}
