package main

import (
	"log"
	"net/http"

	"movieexample.com/rating/internal/controller/rating"
	httpHandler "movieexample.com/rating/internal/handler/http"
	"movieexample.com/rating/internal/repository/memory"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httpHandler.New(ctrl)

	http.HandleFunc("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
