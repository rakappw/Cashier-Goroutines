package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type Product struct {
	Name  string
	Price float64
}

type Customer struct {
	ID         int
	Name       string
	QueueNum   int
	Products   []Product
	TotalPrice float64
}

type Cashier struct {
	ID        int
	Customers int
}

type Store struct {
	Products     []Product
	Customers    []*Customer
	Cashiers     []Cashier
	Transactions []string
	mutex        sync.Mutex
}

func InitializeStore() *Store {
	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(time.Now().UnixNano())

	products := []Product{
		{Name: "Mie Instan", Price: 15000},
		{Name: "Kopi Bubuk", Price: 12500},
		{Name: "Roti", Price: 8000},
		{Name: "Telur", Price: 25000},
		{Name: "Shampoo", Price: 20000},
		{Name: "Gula", Price: 10000},
		{Name: "Minyak Goreng", Price: 30000},
		{Name: "Garam", Price: 5000},
		{Name: "Sabun", Price: 7500},
		{Name: "Pasta gigi", Price: 15000},
	}

	customers := make([]*Customer, 0, 100)
	for i := 0; i < 100; i++ {
		numProducts := rand.Intn(4) + 1
		custProducts := make([]Product, 0, numProducts)
		totalPrice := 0.0

		for j := 0; j < numProducts; j++ {
			product := products[rand.Intn(len(products))]
			custProducts = append(custProducts, product)
			totalPrice += product.Price
		}

		customer := &Customer{
			ID:         i + 1,
			Name:       gofakeit.Name(),
			QueueNum:   i + 1,
			Products:   custProducts,
			TotalPrice: totalPrice,
		}

		customers = append(customers, customer)
	}

	cashiers := make([]Cashier, 5)
	for i := 0; i < 5; i++ {
		cashiers[i] = Cashier{ID: i + 1, Customers: 0}
	}

	return &Store{
		Products:     products,
		Customers:    customers,
		Cashiers:     cashiers,
		Transactions: make([]string, 0),
		mutex:        sync.Mutex{},
	}
}

func (s *Store) ProcessCustomer(cashier *Cashier, customer *Customer, wg *sync.WaitGroup) {
	defer wg.Done()

	processingTime := rand.Intn(4) + 1
	fmt.Printf("Cashier %d serving customer %s (Queue #%d) with %d products. Processing time: %d seconds\n",
		cashier.ID, customer.Name, customer.QueueNum, len(customer.Products), processingTime)

	time.Sleep(time.Duration(processingTime) * time.Second)

	fmt.Printf("\n----- RECEIPT -----\n")
	fmt.Printf("Cashier: %d\n", cashier.ID)
	fmt.Printf("Customer: %s (Queue #%d)\n", customer.Name, customer.QueueNum)
	fmt.Printf("Products:\n")

	for _, product := range customer.Products {
		fmt.Printf("  - %s: Rp%.2f\n", product.Name, product.Price)
	}
	fmt.Printf("Total: Rp%.2f\n", customer.TotalPrice)
	fmt.Printf("-----------------\n\n")

	s.mutex.Lock()
	transaction := fmt.Sprintf("Transaction #%d: Customer %s bought %d products for Rp%.2f",
		len(s.Transactions)+1, customer.Name, len(customer.Products), customer.TotalPrice)
	s.Transactions = append(s.Transactions, transaction)
	cashier.Customers++
	s.mutex.Unlock()
}

func main() {
	fmt.Println("Store Payment Simulation Starting...")
	store := InitializeStore()

	var wg sync.WaitGroup

	cashierWg := sync.WaitGroup{}
	customerChan := make(chan *Customer)

	for i := range store.Cashiers {
		cashierWg.Add(1)
		go func(cashier *Cashier) {
			defer cashierWg.Done()
			for customer := range customerChan {
				wg.Add(1)
				store.ProcessCustomer(cashier, customer, &wg)
			}
		}(&store.Cashiers[i])
	}

	for _, customer := range store.Customers {
		customerChan <- customer
	}

	close(customerChan)
	cashierWg.Wait()
	wg.Wait()

	fmt.Println("\n----- TRANSACTION SUMMARY -----")
	fmt.Printf("Total transactions: %d\n", len(store.Transactions))

	for i, cashier := range store.Cashiers {
		fmt.Printf("Cashier %d processed %d customers\n", i+1, cashier.Customers)
	}

	fmt.Println("\nSample of Transactions:")
	transactionsToShow := 5
	if len(store.Transactions) < transactionsToShow {
		transactionsToShow = len(store.Transactions)
	}

	for i := 0; i < transactionsToShow; i++ {
		fmt.Println(store.Transactions[i])
	}

	fmt.Println("\nStore Payment Simulation Completed!")
}
