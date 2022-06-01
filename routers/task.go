package routers

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/valerianomacuri/task-manager/common"
	"github.com/valerianomacuri/task-manager/controllers"
)

func SetTaskRoutes(router *mux.Router) *mux.Router {
	taskRouter := mux.NewRouter()
	taskRouter.HandleFunc("/tasks", controllers.CreateTask).Methods(http.MethodPost)
	taskRouter.HandleFunc("/tasks/{id}", controllers.UpdateTask).Methods(http.MethodPut)
	taskRouter.HandleFunc("/tasks", controllers.GetTasks).Methods(http.MethodGet)
	taskRouter.HandleFunc("/tasks/{id}", controllers.GetTaskById).Methods(http.MethodGet)
	taskRouter.HandleFunc("/tasks/users/{id}", controllers.GetTasksByUser).Methods(http.MethodGet)
	taskRouter.HandleFunc("/tasks/{id}", controllers.DeleteTask).Methods(http.MethodDelete)
	router.PathPrefix("/tasks").Handler(negroni.New(
		negroni.HandlerFunc(common.Authorize),
		negroni.Wrap(taskRouter),
	))
	return router
}
