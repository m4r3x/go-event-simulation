package main

import (
	"math/rand"
	"time"

	"github.com/msales/pkg/stats"
)

type product string

const (
	bananas product = "bananas"
	apples  product = "apples"
	oranges product = "oranges"
)

type percent int

const (
	visitProduct       percent = 7
	addToCartProduct   percent = 25
	buyProduct         percent = 15
	transactionFailure percent = 1
	minConnections     int     = 300
	maxConnections     int     = 600
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	client, err := stats.NewStatsd("127.0.0.1:8125", "go")
	if err != nil {
		return
	}
	for {
		activeConnections := randInt(minConnections, maxConnections)
		for i := 0; i < activeConnections; i++ {
			go pageVisit(client)
		}
		client.Gauge("active_connections", float64(activeConnections), 1.0, nil)
		time.Sleep(time.Second)
	}
}

func pageVisit(client stats.Stats) {
	client.Inc("page_visit", 1, 1.0, nil)
	if percentageChance(visitProduct) {
		stamp := time.Now()
		time.Sleep(time.Second * time.Duration(randInt(5, 15)))
		productVisit(client, stamp)
	}
}

func productVisit(client stats.Stats, stamp time.Time) {
	product := randomizeProduct()
	client.Inc("product_visit", 1, 1.0, createProductTag(product))
	if percentageChance(addToCartProduct) {
		time.Sleep(time.Second * time.Duration(randInt(5, 15)))
		productAddToCart(client, product, stamp)
	}
}

func productAddToCart(client stats.Stats, product product, stamp time.Time) {
	client.Inc("product_added_to_cart", 1, 1.0, createProductTag(product))
	if percentageChance(buyProduct) {
		time.Sleep(time.Second * time.Duration(randInt(10, 25)))
		productBought(client, product, stamp)
	}
}

func productBought(client stats.Stats, product product, stamp time.Time) {
	client.Inc("product_bought", 1, 1.0, createProductTag(product))
	client.Timing("product_bought_since_visit", time.Since(stamp), 1.0, createProductTag(product))
	if percentageChance(transactionFailure) {
		client.Inc("product_failed_payment", 1, 1.0, createProductTag(product))
	}
}

func randomizeProduct() product {
	if percentageChance(33) {
		return bananas
	}
	if percentageChance(25) {
		return apples
	}

	return oranges
}

func percentageChance(percent percent) bool {
	return randInt(1, 100) < int(percent)
}

func createProductTag(product product) map[string]string {
	tags := make(map[string]string)
	tags["product"] = string(product)

	return tags
}

func randInt(min int, max int) int {
	return min + rand.Intn(max)
}
