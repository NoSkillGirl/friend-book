package main

import (
	"net/http"

	"github.com/NoSkillGirl/friend-book/routers"
)

func main() {
	routers.UserRoutes()
	http.ListenAndServe(":8080", nil)
}
