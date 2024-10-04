package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"client/shopping"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := shopping.NewShoppingServiceClient(conn)

	productsToAdd := []shopping.ProductRequest{
		{Name: "Apples", Quantity: 10},
		{Name: "Bananas", Quantity: 5},
		{Name: "Oranges", Quantity: 7},
	}

	for _, product := range productsToAdd {
		response, err := client.AddProduct(context.Background(), &product)
		if err != nil {
			log.Fatalf("Could not add product: %v", err)
		}
		fmt.Println(response.Message)
	}

	productList, err := client.ListProducts(context.Background(), &shopping.Void{})
	if err != nil {
		log.Fatalf("Could not list products: %v", err)
	}

	fmt.Println("Product List:")
	for _, p := range productList.Products {
		fmt.Printf("Name: %s, Quantity: %d, Purchased: %v\n", p.Name, p.Quantity, p.Purchased)
	}

	updateProduct := &shopping.ProductRequest{Name: "Apples", Quantity: 15}
	updateResponse, err := client.UpdateProduct(context.Background(), updateProduct)
	if err != nil {
		log.Fatalf("Could not update product: %v", err)
	}
	fmt.Println(updateResponse.Message)

	deleteResponse, err := client.DeleteProduct(context.Background(), &shopping.ProductNameRequest{Name: "Bananas"})
	if err != nil {
		log.Fatalf("Could not delete product: %v", err)
	}
	fmt.Println(deleteResponse.Message)

	productList, err = client.ListProducts(context.Background(), &shopping.Void{})
	if err != nil {
		log.Fatalf("Could not list products: %v", err)
	}

	fmt.Println("Updated Product List:")
	for _, p := range productList.Products {
		fmt.Printf("Name: %s, Quantity: %d, Purchased: %v\n", p.Name, p.Quantity, p.Purchased)
	}

	markResponse, err := client.MarkAsPurchased(context.Background(), &shopping.ProductNameRequest{Name: "Oranges"})
	if err != nil {
		log.Fatalf("Could not mark product as purchased: %v", err)
	}
	fmt.Println(markResponse.Message)

	productList, err = client.ListProducts(context.Background(), &shopping.Void{})
	if err != nil {
		log.Fatalf("Could not list products: %v", err)
	}

	fmt.Println("Final Product List:")
	for _, p := range productList.Products {
		fmt.Printf("Name: %s, Quantity: %d, Purchased: %v\n", p.Name, p.Quantity, p.Purchased)
	}
}
