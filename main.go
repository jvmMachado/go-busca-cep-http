package main

import (
	"errors"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", GetCepHandler)
	http.ListenAndServe(":8080", nil)
}

func GetCepHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusBadRequest)
		return
	}

	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		http.Error(w, "cep is empty.", http.StatusBadRequest)
		return
	}

	cep, err := GetCepFromApi(cepParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cep))
}

func GetCepFromApi(cep string) (string, error) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	res, err := http.Get(url)
	if err != nil {
		return "", errors.New("Erro na requisição ao viacep: %v")
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("Erro ao ler o body da requisição: %v")
	}

	return string(data), nil
}