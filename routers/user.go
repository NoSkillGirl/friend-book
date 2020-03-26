package routers

import (
	"github.com/NoSkillGirl/friend-book/controllers"
	"github.com/gorilla/mux"
)

//UserRoutes - routes
func UserRoutes(apiRouter *mux.Router) {

	v1Router := apiRouter.PathPrefix("/v1.0").Subrouter()

	v1Router.Handle("/me", controllers.IsAuthorized(controllers.ShowUser)).Methods("GET")

	v1Router.Handle("/{user_id}", controllers.IsAuthorized(controllers.ShowUser)).Methods("GET")

	v1Router.Handle("/me/update", controllers.IsAuthorized(controllers.UpdateUser)).Methods("PUT")

	v1Router.Handle("/{user_id}/update", controllers.IsAuthorized(controllers.UpdateUser)).Methods("PUT")

	v1Router.Handle("/me/delete", controllers.IsAuthorized(controllers.DeleteUser)).Methods("DELETE")

	v1Router.Handle("/{user_id}/delete", controllers.IsAuthorized(controllers.DeleteUser)).Methods("DELETE")

	v1Router.Handle("/me/friend_request", controllers.IsAuthorized(controllers.FriendRequest)).Methods("POST")

	v1Router.Handle("/me/friend_request", controllers.IsAuthorized(controllers.ActOnFriendRequest)).Methods("PUT")

	v1Router.Handle("/me/friend_requests", controllers.IsAuthorized(controllers.FriendRequests)).Methods("GET")

	// v1Router.Handle("/me/delete_friend_request", controllers.IsAuthorized(controllers.DeleteUser)).Methods("DELETE")

	// v1Router.HandleFunc("/search", controllers.Search).Methods("GET")
}
