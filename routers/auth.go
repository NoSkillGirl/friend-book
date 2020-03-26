package routers

import (
	"github.com/NoSkillGirl/friend-book/controllers"
	"github.com/gorilla/mux"
)

//AuthRoutes - routes for authentication
func AuthRoutes(apiRouter *mux.Router) {

	// Health check routes
	apiRouter.HandleFunc("/health", controllers.HealthCheck).Methods("GET")

	// auth routes
	apiRouter.HandleFunc("/signup", controllers.Signup).Methods("POST")
	apiRouter.HandleFunc("/login", controllers.Login).Methods("POST")
}
