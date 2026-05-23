package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverURL     = "http://localhost:8080/cotacao"
	clientTimeout = 300 * time.Millisecond
	outputFile    = "cotacao.txt"
)

type cotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("timeout do cliente: %v", err)
		}
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("server returned %d: %s", res.StatusCode, string(body))
	}

	var data cotacaoResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	content := fmt.Sprintf("Dólar: %s", data.Bid)
	if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
}
