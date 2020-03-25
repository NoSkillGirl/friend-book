package routers

import (
	"net/http"

	"github.com/NoSkillGirl/friend-book/controllers"
	"github.com/gorilla/mux"
)

//UserRoutes - routes
func UserRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1.0/", controllers.HealthCheck)
	r.HandleFunc("/api/v1.0/signup", controllers.UserSignup)
	r.HandleFunc("/api/v1.0/login", controllers.UserLogin)
	r.HandleFunc("/api/v1.0/users/{user_id}/update", controllers.UpdateUser)
	r.HandleFunc("/api/v1.0/users/{user_id}/delete", controllers.DeleteUser)
	r.HandleFunc("/api/v1.0/users/{user_id}", controllers.ShowUser)
	r.HandleFunc("/api/v1.0/users/{user_id}/friend_request", controllers.FriendRequest)
	r.HandleFunc("/api/v1.0/users/{user_id}/search", controllers.Search)
	// r.HandleFunc("/api/v1.0/users/{user_id}/friends/{friend_id}/delete", controllers.DeleteFriend)
	r.HandleFunc("/api/v1.0/users/{user_id}/friend_requests", controllers.FriendRequests)
	// r.HandleFunc("/api/v1.0/users/{user_id}/friends/{friend_id}/response", controllers.FriendRequestResponse)
	http.Handle("/", r)
}
