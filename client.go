package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	println("Inicio do processo...\n")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	cambio, err := GetCambio(ctx)
	if err != nil {
		println("\nErro ao request cambio")
	}

	println("Dólar: " + string(cambio))
	println("\nfim do processo")
}

func GetCambio(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
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
	return body, nil
}
