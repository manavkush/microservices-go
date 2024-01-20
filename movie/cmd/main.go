package main

import (
	"log"
	"net/http"

	"movieexample.com/movie/internal/controller/movie"
	metadataGateway "movieexample.com/movie/internal/gateway/metadata/http"
	ratingGateway "movieexample.com/movie/internal/gateway/rating/http"
	httpHandler "movieexample.com/movie/internal/handler/http"
)

func main() {
	log.Println("Starting the movie service")
	metadataGateway := metadataGateway.New("localhost:8081")
	ratingGateway := ratingGateway.New("localhost:8082")

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httpHandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
