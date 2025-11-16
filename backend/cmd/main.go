package main

import (
	"log"
	"moment/internal/handler"
	"net/http"
)

func main() {
	h := handler.NewHandler()

	mux := http.NewServeMux()

	// 所有后端路由统一加 /api
	mux.Handle("/api/", http.StripPrefix("/api", h))

	log.Println("listen on 127.0.0.1:8080")
	http.ListenAndServe("127.0.0.1:8080", mux)
}
