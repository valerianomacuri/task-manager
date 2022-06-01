package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/valerianomacuri/task-manager/controllers"
)

func SetUserRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/users/register", controllers.Register).Methods(http.MethodPost)
	router.HandleFunc("/users/login", controllers.Login).Methods(http.MethodPost)
	return router
}
