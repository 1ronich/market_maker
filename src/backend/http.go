package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
func messageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("salam")
	enableCors(w)

	if r.Method != http.MethodPost {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received message: %s\n", string(body))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message received"))
}

func maissan() {
	http.HandleFunc("/message", messageHandler)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
