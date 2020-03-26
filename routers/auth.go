package routers

import (
	"github.com/NoSkillGirl/friend-book/controllers"
	"github.com/gorilla/mux"
)

//AuthRoutes - routes for authentication
func AuthRoutes(router *mux.Router, apiRouter *mux.Router) {

	// Health check routes
	apiRouter.HandleFunc("/health", controllers.HealthCheck).Methods("GET")

	// Html Pages
	router.HandleFunc("/login", controllers.GetLoginPage).Methods("GET")
	router.HandleFunc("/register", controllers.GetSignUpPage).Methods("GET")

	// auth routes
	apiRouter.HandleFunc("/register", controllers.Signup).Methods("POST")
	apiRouter.HandleFunc("/login", controllers.Login).Methods("POST")

}
