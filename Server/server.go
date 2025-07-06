package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cambio struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", HttpHandler)
	http.ListenAndServe(":8080", mux)
}

func HttpHandler(w http.ResponseWriter, r *http.Request) {
	println("Procesando request...")
	defer println("request processado...")
	cambio, err := GetLastCambio("USD-BRL")
	if err != nil {
		println("\nErro ao request cambio")
	}
	res, err := json.Marshal(cambio.USDBRL)
	if err != nil {
		println(err)
	}

	SalvaCotacao(cambio.USDBRL.Bid, r)
	CriarArquivo("cotacao.txt", "Dólar: "+string(res))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cambio.USDBRL.Bid))
}

func GetLastCambio(moeda string) (*Cambio, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/"+moeda, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		log.Println("Erro request externa:", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, error := io.ReadAll(res.Body)
	if error != nil {
		return nil, error
	}
	var cambio Cambio
	error = json.Unmarshal(body, &cambio)
	if error != nil {
		return nil, error
	}
	return &cambio, nil
}

func SalvaCotacao(bid string, r *http.Request) {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		log.Fatal("Erro ao abrir banco:", err)
	}
	defer db.Close()

	ctxDB, cancelDB := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancelDB()

	_, err = db.ExecContext(ctxDB, "INSERT INTO cotacoes (bid) VALUES (?)", bid)
	if err != nil {
		log.Println("Erro ao salvar no banco:", err)
	}
}

func CriarArquivo(arquivoNome string, content string) {
	file, err := os.Create(arquivoNome)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	tamanho, err := file.Write([]byte(content))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Arquivo criado com sucesso! Tamanho %d bytes", tamanho)
}
