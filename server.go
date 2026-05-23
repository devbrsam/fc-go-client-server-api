package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	apiURL         = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiTimeout     = 200 * time.Millisecond
	dbTimeout      = 10 * time.Millisecond
	serverAddr     = ":8080"
	cotacaoPath    = "/cotacao"
	dbFile         = "cotacoes.db"
)

type usdBrlAPI struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type cotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc(cotacaoPath, func(w http.ResponseWriter, r *http.Request) {
		handleCotacao(w, r, db)
	})

	log.Printf("server listening on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS cotacoes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bid TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func handleCotacao(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	bid, err := fetchBid(r.Context())
	if err != nil {
		if isTimeout(err) {
			log.Printf("timeout da API: %v", err)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := saveBid(r.Context(), db, bid); err != nil {
		if isTimeout(err) {
			log.Printf("timeout do banco: %v", err)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacaoResponse{Bid: bid})
}

func fetchBid(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var data usdBrlAPI
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	return data.USDBRL.Bid, nil
}

func saveBid(ctx context.Context, db *sql.DB, bid string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)", bid)
	return err
}

func isTimeout(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}
