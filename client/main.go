package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "client/shopping" // Замените на путь к сгенерированным файлам
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewShoppingServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nВыберите операцию:")
		fmt.Println("1: Добавить продукт")
		fmt.Println("2: Показать продукт")
		fmt.Println("3: Обновить продукт")
		fmt.Println("4: Удалить продукт")
		fmt.Println("5: Показать список покупок")
		fmt.Println("6: Отметить продукт как купленный")
		fmt.Println("0: Выйти")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, _ := strconv.Atoi(input)

		switch choice {
		case 1:
			addProduct(client, reader)
		case 2:
			getProduct(client, reader)
		case 3:
			updateProduct(client, reader)
		case 4:
			deleteProduct(client, reader)
		case 5:
			listProducts(client)
		case 6:
			markAsPurchased(client, reader)
		case 0:
			return
		default:
			fmt.Println("Некорректный выбор. Попробуйте снова.")
		}
	}
}

func addProduct(client pb.ShoppingServiceClient, reader *bufio.Reader) {
	fmt.Print("Введите название продукта: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Введите количество: ")
	qtyStr, _ := reader.ReadString('\n')
	qty, _ := strconv.Atoi(strings.TrimSpace(qtyStr))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.AddProduct(ctx, &pb.ProductRequest{Name: name, Quantity: int32(qty)})
	if err != nil {
		log.Fatalf("Ошибка при добавлении продукта: %v", err)
	}
	fmt.Printf("Продукт добавлен: %s\n", res.Message)
}

func getProduct(client pb.ShoppingServiceClient, reader *bufio.Reader) {
	fmt.Print("Введите название продукта для получения: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetProduct(ctx, &pb.ProductNameRequest{Name: name})
	if err != nil {
		log.Fatalf("Ошибка при получении продукта: %v", err)
	}
	fmt.Printf("Продукт: %s, Количество: %d, Куплен: %t\n", res.Name, res.Quantity, res.Purchased)
}

func updateProduct(client pb.ShoppingServiceClient, reader *bufio.Reader) {
	fmt.Print("Введите название продукта для обновления: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Введите новое количество: ")
	qtyStr, _ := reader.ReadString('\n')
	qty, _ := strconv.Atoi(strings.TrimSpace(qtyStr))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.UpdateProduct(ctx, &pb.ProductRequest{Name: name, Quantity: int32(qty)})
	if err != nil {
		log.Fatalf("Ошибка при обновлении продукта: %v", err)
	}
	fmt.Printf("Обновление успешно: %s\n", res.Message)
}

func deleteProduct(client pb.ShoppingServiceClient, reader *bufio.Reader) {
	fmt.Print("Введите название продукта для удаления: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.DeleteProduct(ctx, &pb.ProductNameRequest{Name: name})
	if err != nil {
		log.Fatalf("Ошибка при удалении продукта: %v", err)
	}
	fmt.Printf("Удаление успешно: %s\n", res.Message)
}

func listProducts(client pb.ShoppingServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.ListProducts(ctx, &pb.Void{})
	if err != nil {
		log.Fatalf("Ошибка при выводе списка продуктов: %v", err)
	}

	fmt.Println("Список продуктов:")
	for _, item := range res.Products {
		fmt.Printf("Название: %s, Количество: %d, Куплен: %t\n", item.Name, item.Quantity, item.Purchased)
	}
}

func markAsPurchased(client pb.ShoppingServiceClient, reader *bufio.Reader) {
	fmt.Print("Введите название продукта для отметки как купленный: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.MarkAsPurchased(ctx, &pb.ProductNameRequest{Name: name})
	if err != nil {
		log.Fatalf("Ошибка при отметке продукта как купленного: %v", err)
	}
	fmt.Printf("Продукт успешно отмечен как купленный: %s\n", res.Message)
}
