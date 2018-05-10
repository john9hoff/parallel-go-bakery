package main

import (
	"fmt"
	"os"
	"math/rand"
	"time"
	"strconv"
)

type Order struct {
	num int
	sender chan int
}

func fib(n int)(int) {
	switch n {
	case 0, 1:
		return n
	default:
		return fib(n-1) + fib(n-2)
	}
}

func worker(c chan Order) {
	for {
		order := <-c // worker takes an order
		order.sender <- fib(order.num) // worker returns the processed order to customer
	}
}

func customer(c chan Order, t chan bool) {

	sleepTime := time.Duration(rand.Intn(10000)) * time.Millisecond

	// sleep for a random duration
	time.Sleep(sleepTime)

	resultChan := make(chan int)

	// make an order from a random fib of order number
	myOrder := Order{sender: resultChan, num: rand.Intn(43)}

	fmt.Println("The customer was sleeping for", sleepTime)
	fmt.Println("They received order number", myOrder.num)

	c <- myOrder // put order in channel for worker

	result := <- resultChan // get result from worker

		fmt.Println("They were served a donut weighing", result, "pounds")

	t <- true // put confirmation in finished customers channel

}

func main() {

	startTime := time.Now()
	args := os.Args[1:]
	numWorkers, err1 := strconv.Atoi(args[0])
	numCustomers, err2 := strconv.Atoi(args[1])

	if err1 != nil {
		fmt.Println(err1)
		os.Exit(2)
	}
	if err2 != nil {
		fmt.Println(err2)
		os.Exit(2)
	}

	// create channel for orders from customers
	workerLine := make(chan Order, numCustomers)

	// create channel to put customers in when they're served
	finishedChannel := make(chan bool, numCustomers)

	// create workers
	for i := 0; i < numWorkers; i++ {
		go worker(workerLine)
		fmt.Println("Worker", i, "is ready to work.")
	}

	// create customers
	for i := 0; i < numCustomers; i++ {
		go customer(workerLine, finishedChannel)
		fmt.Println("Customer", i, "is ready to order.")
	}

	// wait for all customers to be served.
	for i := 0; i < numCustomers; i++ {
		<- finishedChannel
		fmt.Println("Customer", i, "was served")
	}

	timeElapsed := time.Since(startTime)

	fmt.Println("The bakery served everyone in", timeElapsed)
	fmt.Println("The bakery has closed.")

}
