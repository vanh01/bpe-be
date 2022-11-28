package main

import "bpe/router"

func main() {
	server := router.NewServer()
	server.Run()
}
