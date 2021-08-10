package order_server

import (
	"context"
	"fmt"
	_ "github.com/google/uuid"
	"github.com/mertturkman/ordergrpc/orderpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	_ "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var collection *mongo.Collection

type server struct {
}

type orderItem struct {
	Id              primitive.ObjectID  `bson:"_id,omitempty"`
	OrderNumber     string              `bson:"orderNumber"`
	CreatedDateTime time.Time           `bson:"createdDateTime"`
	User            userItem            `bson:"user"`
	ShippingAddress shippingAddressItem `bson:"shippingAddress"`
	Lines           []orderLineItem     `bson:"lines"`
	TotalAmount     amountItem          `bson:"totalAmount"`
}

type userItem struct {
	Id          string
	Username    string
	Mail        string
	PhoneNumber string
}

type shippingAddressItem struct {
	Id      string
	City    string
	Country string
	County  string
	Detail  string
}

type orderLineItem struct {
	Name     string
	Sku      string
	Quantity int
	Price    amountItem
}

type amountItem struct {
	Value    float64
	Currency string
}

func (s server) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	fmt.Println("Create order started.")
	order := request.GetOrder()

	data := orderItem{
		CreatedDateTime: time.Now(),
		OrderNumber:     order.OrderNumber,
		User: userItem{
			Id:          order.User.Id,
			Username:    order.User.UserName,
			Mail:        order.User.Mail,
			PhoneNumber: order.User.PhoneNumber,
		},
		ShippingAddress: shippingAddressItem{
			Id:      order.ShippingAddress.Id,
			Country: order.ShippingAddress.Country,
			City:    order.ShippingAddress.City,
			County:  order.ShippingAddress.County,
			Detail:  order.ShippingAddress.Detail,
		},
		TotalAmount: amountItem{
			Value: order.TotalAmount.Value,
		},
	}

	for _, lineItem := range order.Items {
		data.Lines = append(data.Lines, orderLineItem{
			Name:     lineItem.Name,
			Sku:      lineItem.Sku,
			Quantity: int(lineItem.Quantity),
			Price: amountItem{
				Value:    lineItem.Price.Value,
				Currency: lineItem.Price.Currency,
			},
		})
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Collection insert error: %v", err),
		)
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}

	return &orderpb.CreateOrderResponse{
		Id: id.Hex(),
	}, nil

}

func (s server) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	fmt.Println("Read order started.")

	orderID := request.GetId()

	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	data := &orderItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot order blog with specified ID: %v", err),
		)
	}

	orderItem := &orderpb.Order{
		Id:              data.Id.Hex(),
		CreatedDateTime: timestamppb.New(data.CreatedDateTime),
		OrderNumber:     data.OrderNumber,
		User: &orderpb.User{
			Id:          data.User.Id,
			UserName:    data.User.Username,
			Mail:        data.User.Mail,
			PhoneNumber: data.User.PhoneNumber,
		},
		TotalAmount: &orderpb.Amount{
			Currency: data.TotalAmount.Currency,
			Value:    data.TotalAmount.Value,
		},
		ShippingAddress: &orderpb.ShippingAddress{
			Id:      data.ShippingAddress.Id,
			Country: data.ShippingAddress.Country,
			City:    data.ShippingAddress.City,
			County:  data.ShippingAddress.County,
			Detail:  data.ShippingAddress.Detail,
		},
	}

	for _, lineItem := range data.Lines {
		orderItem.Items = append(orderItem.Items, &orderpb.OrderLine{
			Sku:      lineItem.Sku,
			Quantity: int32(lineItem.Quantity),
			Name:     lineItem.Name,
			Price: &orderpb.Amount{
				Currency: lineItem.Price.Currency,
				Value:    lineItem.Price.Value,
			},
		})
	}

	return &orderpb.GetOrderResponse{
		Order: orderItem,
	}, nil
}

func (s server) UpdateOrder(ctx context.Context, request *orderpb.UpdateOrderRequest) (*emptypb.Empty, error) {
	fmt.Println("Update order started.")
	order := request.GetOrder()

	oid, err := primitive.ObjectIDFromHex(order.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	data := orderItem{
		Id:              oid,
		CreatedDateTime: time.Now(),
		OrderNumber:     order.OrderNumber,
		User: userItem{
			Id:          order.User.Id,
			Username:    order.User.UserName,
			Mail:        order.User.Mail,
			PhoneNumber: order.User.PhoneNumber,
		},
		ShippingAddress: shippingAddressItem{
			Id:      order.ShippingAddress.Id,
			Country: order.ShippingAddress.Country,
			City:    order.ShippingAddress.City,
			County:  order.ShippingAddress.County,
			Detail:  order.ShippingAddress.Detail,
		},
		TotalAmount: amountItem{
			Value: order.TotalAmount.Value,
		},
	}

	for _, lineItem := range order.Items {
		data.Lines = append(data.Lines, orderLineItem{
			Name:     lineItem.Name,
			Sku:      lineItem.Sku,
			Quantity: int(lineItem.Quantity),
			Price: amountItem{
				Value:    lineItem.Price.Value,
				Currency: lineItem.Price.Currency,
			},
		})
	}

	oid, err = primitive.ObjectIDFromHex(order.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": data}
	res, err := collection.UpdateOne(context.Background(), filter, update)

	if res.ModifiedCount <= 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot updated order with specified ID: %v", err),
		)
	}

	return &emptypb.Empty{}, nil
}

func (s server) DeleteOrder(ctx context.Context, request *orderpb.DeleteOrderRequest) (*emptypb.Empty, error) {
	orderID := request.GetId()

	oid, err := primitive.ObjectIDFromHex(orderID)

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID."),
		)
	}

	filter := bson.M{"_id": oid}
	res, err := collection.DeleteOne(context.Background(), filter)
	if res.DeletedCount <= 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot deleted order with specified ID: %v", err),
		)
	}

	return &emptypb.Empty{}, nil
}
