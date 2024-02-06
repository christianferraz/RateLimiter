package main

import (
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRateLimiterUnderLimitWithToken(t *testing.T) {
	totalRequests := 9 // Número de requisições menor que o limite
	var wg sync.WaitGroup
	wg.Add(totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func(i int) {
			defer wg.Done()
			req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
			if err != nil {
				t.Errorf("Erro ao criar requisição %d: %v", i, err)
				return
			}
			req.Header.Set("API_KEY", "token_1")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Erro ao fazer requisição %d: %v", i, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Requisição %d falhou ou foi inesperadamente limitada. Status: %v", i, resp.StatusCode)
				os.Exit(0)
			}
		}(i)
	}

	wg.Wait()
	t.Logf("Todas as requisições foram feitas com sucesso")
}

func TestRateLimiterUnderLimit(t *testing.T) {
	totalRequests := 9 // Número de requisições menor que o limite
	var wg sync.WaitGroup
	wg.Add(totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func(i int) {
			defer wg.Done()
			resp, err := http.Get("http://localhost:8080/")

			if err != nil || resp.StatusCode != http.StatusOK {
				t.Errorf("Requisição falhou ou foi inesperadamente limitada %v", resp.StatusCode)
			}
			defer resp.Body.Close()
		}(i)
	}

	wg.Wait()
	t.Logf("Todas as requisições foram feitas com sucesso")
}

func TestRateLimiterOverLimit(t *testing.T) {
	totalRequests := 15
	var blockedRequests int
	var wg sync.WaitGroup
	results := make(chan bool, totalRequests)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get("http://localhost:8080/")
			time.Sleep(1000 * time.Millisecond)
			if err != nil {
				results <- false // Indica falha na requisição, não necessariamente rate limit
				return
			}
			defer resp.Body.Close()

			results <- resp.StatusCode == http.StatusTooManyRequests
		}()
	}

	wg.Wait()
	close(results)

	for result := range results {
		if result {
			blockedRequests++
		}
	}

	// Verifica se alguma requisição foi bloqueada pelo rate limiter
	if blockedRequests == 0 {
		t.Errorf("O rate limiter não bloqueou nenhuma requisição")
	}
}
