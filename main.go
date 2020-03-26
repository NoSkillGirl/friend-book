package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NoSkillGirl/friend-book/routers"
	"github.com/gorilla/mux"
)

const (
	STATIC_DIR = "/assets/"
)

func main() {
	router := mux.NewRouter()
	router.
		PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	apiRouter := router.PathPrefix("/api").Subrouter()

	routers.AuthRoutes(router, apiRouter)
	routers.UserRoutes(apiRouter)

	http.Handle("/", router)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("server is starting up on http://127.0.0.1:8000 .........")
	log.Fatal(srv.ListenAndServe())
}
