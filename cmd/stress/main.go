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
	numRequests = 10000
	concurrency = 50
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
	fmt.Println("Initializing ID pool for realistic requests...")
	pool := fetchIDPool()
	if len(pool.CategoryIDs) == 0 || len(pool.ProductIDs) == 0 {
		fmt.Println("Error: Could not fetch IDs for testing. Ensure server has data.")
		return
	}

	fmt.Printf("Starting enhanced stress test: %d requests with %d concurrency\n", numRequests, concurrency)
	
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
	fmt.Printf("\n--- Enhanced Stress Test Results ---\n")
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
	resp, _ := http.Get(baseURL + "/category?limit=50")
	var catData struct {
		Data []struct{ ID string } `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&catData)
	for _, c := range catData.Data {
		pool.CategoryIDs = append(pool.CategoryIDs, c.ID)
	}
	resp.Body.Close()

	// Fetch Products
	resp, _ = http.Get(baseURL + "/products/?limit=50")
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
	case r < 0.3: // 30% List Products
		return doRequest("GET", "/products/?limit=10")
	case r < 0.5: // 20% Get Product By ID
		id := pool.ProductIDs[rand.Intn(len(pool.ProductIDs))]
		return doRequest("GET", "/products/"+id)
	case r < 0.65: // 15% Get Product By Category
		id := pool.CategoryIDs[rand.Intn(len(pool.CategoryIDs))]
		return doRequest("GET", "/products/category/"+id)
	case r < 0.8: // 15% List Categories
		return doRequest("GET", "/category?limit=10")
	case r < 0.9: // 10% Create Product
		return doPOSTProduct(pool.CategoryIDs[0])
	default: // 10% Other (Mocked PUT/DELETE)
		return Result{"OTHER (PUT/DELETE)", 200, nil}
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

func doPOSTProduct(catID string) Result {
	body, _ := json.Marshal(map[string]interface{}{
		"name":         fmt.Sprintf("Stress-%d", rand.Int()),
		"slug":         fmt.Sprintf("s-%d-%d", time.Now().UnixNano(), rand.Int()),
		"Description": "stress",
		"category_id":  catID,
		"price":        1.99,
	})
	resp, err := http.Post(baseURL+"/products/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return Result{"POST /products/", 0, err}
	}
	defer resp.Body.Close()
	return Result{"POST /products/", resp.StatusCode, nil}
}
