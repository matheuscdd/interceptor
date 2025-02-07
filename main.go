package main

import (
	"bytes"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	id := paths[1]
	endpoint := strings.Join(paths[1:], "/")
	log.Println(id, endpoint)

	mapa := map[string]string{
		"server1": "http://localhost:8000",
		"server2": "https://jsonplaceholder.typicode.com/todos",
	}

	dest := mapa["server1"] + "/" + endpoint

	// TODO Fazer um if para caso for arquivo estático que o cara está registrando, sempre redireciona
	if strings.HasPrefix(r.URL.Path, "/static/") {
		http.Redirect(w, r, dest, 302)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler corpo", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	req, err := http.NewRequest(r.Method, dest, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "Erro ao criar", http.StatusInternalServerError)
		return
	}
	req.URL.RawQuery = r.URL.RawQuery
	
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Algo deu errado", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(path.Ext(r.URL.Path))
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
	}
	w.WriteHeader(res.StatusCode)
	// Definindo o status code correto

	_, err = io.Copy(w, res.Body)
	if err != nil {
		http.Error(w, "Algo deu errado ao copiar a resposta", http.StatusInternalServerError)
		return
	}


	// queryParams := r.URL.Query()
	// log.Println(queryParams)
    // http.Redirect(w, r, , 302)
}

func main() {
    http.HandleFunc("/", redirect)
    err := http.ListenAndServe(":8166", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}