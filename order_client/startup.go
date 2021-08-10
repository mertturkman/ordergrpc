package order_client

import (
	"context"
	_ "context"
	"fmt"
	"github.com/mertturkman/ordergrpc/orderpb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	_ "io"
	"log"
	"time"
)

func Init() {
	fmt.Println("Starting Order Client.")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50052", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := orderpb.NewOrderServiceClient(cc)

	fmt.Println("Creating order.")

	order := &orderpb.Order{
		OrderNumber:     "33321342",
		CreatedDateTime: timestamppb.New(time.Now()),
		User: &orderpb.User{
			Id:          "eb11400c-5c97-4096-a6be-434714d7b8e4",
			UserName:    "mertturkman",
			Mail:        "info@mertturkman.com",
			PhoneNumber: "5382900657",
		},
		ShippingAddress: &orderpb.ShippingAddress{
			Country: "TURKEY",
			City:    "Istanbul",
			County:  "Maltepe",
			Detail:  "Altıntepe",
		},
		Items: []*orderpb.OrderLine{
			{
				Quantity: 2,
				Name:     "Computer",
				Sku:      "FF1234HVG242",
				Price: &orderpb.Amount{
					Value:    22.66,
					Currency: "945",
				},
			},
			{
				Quantity: 1,
				Name:     "Mouse",
				Sku:      "MM190HVG004",
				Price: &orderpb.Amount{
					Value:    2.21,
					Currency: "945",
				},
			},
		},
		TotalAmount: &orderpb.Amount{
			Value:    24.87,
			Currency: "994",
		},
	}

	createOrderRes, err := c.CreateOrder(context.Background(), &orderpb.CreateOrderRequest{Order: order})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Order created: %v", createOrderRes)
	orderID := createOrderRes.Id

	readOrderReq := &orderpb.GetOrderRequest{Id: orderID}
	readOrderRes, readOrderErr := c.GetOrder(context.Background(), readOrderReq)
	if readOrderErr != nil {
		fmt.Printf("Error occured while reading: %v \n", readOrderErr)
	}
	fmt.Printf("Order was read: %v \n", readOrderRes)

	readOrderRes.Order.Items[0].Quantity = 1
	readOrderRes.Order.Items[0].Price.Value = 11.33
	readOrderRes.Order.TotalAmount.Value = 13.54

	_, updateErr := c.UpdateOrder(context.Background(), &orderpb.UpdateOrderRequest{Order: readOrderRes.Order})
	if updateErr != nil {
		fmt.Printf("Error occured while updating: %v \n", updateErr)
	}
	fmt.Printf("Order was updated.")

	_, deleteErr := c.DeleteOrder(context.Background(), &orderpb.DeleteOrderRequest{Id: orderID})
	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Order was deleted.")
}
