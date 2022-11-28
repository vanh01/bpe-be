package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() {
	fmt.Println("http Server is running on http://localhost:8080")
	r := mux.NewRouter()

	MapRouters(r)

	log.Fatal(http.ListenAndServe(":8080", r))
}
