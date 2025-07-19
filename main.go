package main

import "sub/infraestructure/adapters"

func main() {
	rabbit := adapters.NewConn()
	rabbit.ListenToQueue()
}