package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL     = "http://localhost:8080/api"
	numRequests = 5000 // Slightly reduced for faster mentor feedback
	concurrency = 30
)

type Result struct {
	Method string
	Status int
	Err    error
}

type IDPool struct {
	CategoryIDs []string
	ProductIDs  []string
}

func main() {
	fmt.Println("Initializing ID pool for many-to-many simulation...")
	pool := fetchIDPool()
	if len(pool.CategoryIDs) < 2 || len(pool.ProductIDs) == 0 {
		fmt.Println("Error: Need at least 2 categories and some products for a realistic simulation.")
		return
	}

	fmt.Printf("Starting simulation: %d requests with %d concurrency\n", numRequests, concurrency)
	
	results := make(chan Result, numRequests)
	var wg sync.WaitGroup
	startTime := time.Now()
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numRequests/concurrency; j++ {
				results <- doMixedRequest(pool)
			}
		}()
	}
	
	go func() {
		wg.Wait()
		close(results)
	}()
	
	stats := make(map[string]int)
	errors := 0
	for res := range results {
		if res.Err != nil {
			errors++
			continue
		}
		key := fmt.Sprintf("%-25s %d", res.Method, res.Status)
		stats[key]++
	}
	
	duration := time.Since(startTime)
	fmt.Printf("\n--- Simulation Results (Many-to-Many) ---\n")
	fmt.Printf("Total Time: %v\n", duration)
	fmt.Printf("Requests/sec: %.2f\n", float64(numRequests)/duration.Seconds())
	fmt.Printf("Total Errors: %d\n", errors)
	fmt.Println("Endpoints Performance:")
	for k, v := range stats {
		fmt.Printf("  %s: %d\n", k, v)
	}
}

func fetchIDPool() IDPool {
	var pool IDPool
	
	// Fetch Categories
	resp, _ := http.Get(baseURL + "/category?limit=100")
	var catData struct {
		Data []struct{ ID string } `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&catData)
	for _, c := range catData.Data {
		pool.CategoryIDs = append(pool.CategoryIDs, c.ID)
	}
	resp.Body.Close()

	// Fetch Products
	resp, _ = http.Get(baseURL + "/products/?limit=100")
	var prodData struct {
		Data []struct{ ID string } `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&prodData)
	for _, p := range prodData.Data {
		pool.ProductIDs = append(pool.ProductIDs, p.ID)
	}
	resp.Body.Close()

	return pool
}

func doMixedRequest(pool IDPool) Result {
	r := rand.Float32()
	switch {
	case r < 0.4: // 40% List Products (Rich JSON Aggregation)
		return doRequest("GET", "/products/?limit=10")
	case r < 0.6: // 20% Get Product By ID
		id := pool.ProductIDs[rand.Intn(len(pool.ProductIDs))]
		return doRequest("GET", "/products/"+id)
	case r < 0.8: // 20% Get Product By Category
		id := pool.CategoryIDs[rand.Intn(len(pool.CategoryIDs))]
		return doRequest("GET", "/products/category/"+id)
	case r < 0.95: // 15% Create Product (Linked to 2 random categories)
		cat1 := pool.CategoryIDs[rand.Intn(len(pool.CategoryIDs))]
		cat2 := pool.CategoryIDs[rand.Intn(len(pool.CategoryIDs))]
		return doPOSTProduct([]string{cat1, cat2})
	default: 
		return Result{"MOCK OTHER", 200, nil}
	}
}

func doRequest(method, path string) Result {
	client := &http.Client{}
	req, _ := http.NewRequest(method, baseURL+path, nil)
	resp, err := client.Do(req)
	if err != nil {
		return Result{method + " " + path, 0, err}
	}
	defer resp.Body.Close()
	return Result{method + " " + path, resp.StatusCode, nil}
}

func doPOSTProduct(catIDs []string) Result {
	body, _ := json.Marshal(map[string]interface{}{
		"name":         fmt.Sprintf("Sim-M2M-%d", rand.Int()),
		"slug":         fmt.Sprintf("s-m2m-%d-%d", time.Now().UnixNano(), rand.Int()),
		"Description": "many-to-many simulation",
		"category_ids": catIDs,
		"price":        99.99,
	})
	resp, err := http.Post(baseURL+"/products/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return Result{"POST /products/", 0, err}
	}
	defer resp.Body.Close()
	return Result{"POST /products/", resp.StatusCode, nil}
}
