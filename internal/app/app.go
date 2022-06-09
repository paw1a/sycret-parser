package app

import (
	"fmt"
	"github.com/paw1a/sycret-parser/internal/handler"
	"log"
	"net/http"
)

const (
	Port = 8080
)

func Run() {
	http.HandleFunc("/api/doc", handler.DocEndpoint)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))
}
