package main

import (
	"log"
	"moment/internal/handler"
	"net/http"
)

func main() {
	h := handler.NewHandler()

	log.Println("listen on 127.0.0.1:8080") // 可选：修改提示
	http.ListenAndServe("127.0.0.1:8080", h)
}
