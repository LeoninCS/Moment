package main

import (
	"log"
	"moment/internal/handler"
	"net/http"
)

func main() {
	h := handler.NewHandler()

	log.Println("listen on :8080")
	http.ListenAndServe(":8080", h)
}
