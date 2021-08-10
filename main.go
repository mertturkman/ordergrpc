package main

import (
	"github.com/mertturkman/ordergrpc/order_client"
	"github.com/mertturkman/ordergrpc/order_server"
)

func main() {
	go order_server.Init()
	order_client.Init()
}
