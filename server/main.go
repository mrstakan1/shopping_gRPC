package main

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"server/shopping"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type ShoppingServiceServer struct {
	shopping.UnimplementedShoppingServiceServer
	redisClient *redis.Client
}

func NewShoppingServiceServer() *ShoppingServiceServer {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	return &ShoppingServiceServer{
		redisClient: rdb,
	}
}

func (s *ShoppingServiceServer) AddProduct(ctx context.Context, req *shopping.ProductRequest) (*shopping.ProductResponse, error) {
	key := "product:" + req.Name
	err := s.redisClient.HSet(ctx, key, "quantity", req.Quantity, "purchased", false).Err()
	if err != nil {
		return nil, err
	}

	return &shopping.ProductResponse{Message: "Product added"}, nil
}

func (s *ShoppingServiceServer) GetProduct(ctx context.Context, req *shopping.ProductNameRequest) (*shopping.Product, error) {
	key := "product:" + req.Name
	data, err := s.redisClient.HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("product not found")
	}

	quantity, _ := strconv.Atoi(data["quantity"])
	purchased, _ := strconv.ParseBool(data["purchased"])

	return &shopping.Product{
		Name:      req.Name,
		Quantity:  int32(quantity),
		Purchased: purchased,
	}, nil
}

func (s *ShoppingServiceServer) UpdateProduct(ctx context.Context, req *shopping.ProductRequest) (*shopping.ProductResponse, error) {
	key := "product:" + req.Name
	err := s.redisClient.HSet(ctx, key, "quantity", req.Quantity).Err()
	if err != nil {
		return nil, err
	}

	return &shopping.ProductResponse{Message: "Product updated"}, nil
}

func (s *ShoppingServiceServer) DeleteProduct(ctx context.Context, req *shopping.ProductNameRequest) (*shopping.ProductResponse, error) {
	key := "product:" + req.Name
	err := s.redisClient.Del(ctx, key).Err()
	if err != nil {
		return nil, err
	}

	return &shopping.ProductResponse{Message: "Product deleted"}, nil
}

func (s *ShoppingServiceServer) ListProducts(ctx context.Context, req *shopping.Void) (*shopping.ProductList, error) {
	keys, err := s.redisClient.Keys(ctx, "product:*").Result()
	if err != nil {
		return nil, err
	}

	var products []*shopping.Product
	for _, key := range keys {
		name := key[8:]
		data, _ := s.redisClient.HGetAll(ctx, key).Result()

		quantity, _ := strconv.Atoi(data["quantity"])
		purchased, _ := strconv.ParseBool(data["purchased"])

		products = append(products, &shopping.Product{
			Name:      name,
			Quantity:  int32(quantity),
			Purchased: purchased,
		})
	}

	return &shopping.ProductList{Products: products}, nil
}

func (s *ShoppingServiceServer) MarkAsPurchased(ctx context.Context, req *shopping.ProductNameRequest) (*shopping.ProductResponse, error) {
	key := "product:" + req.Name
	err := s.redisClient.HSet(ctx, key, "purchased", true).Err()
	if err != nil {
		return nil, err
	}

	return &shopping.ProductResponse{Message: "Product marked as purchased"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	shopping.RegisterShoppingServiceServer(grpcServer, NewShoppingServiceServer())

	fmt.Println("Server running on port 9090")

	if err := grpcServer.Serve(lis); err != nil {
		grpcServer.Stop()
		panic(err)
	}
}
