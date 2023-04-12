package main

import (
	"log"
	"net/http"
	"soap-example/server"
)

func main() {
	service := &server.CalculatorService{}

	http.HandleFunc("/soap", func(w http.ResponseWriter, r *http.Request) {
		server.SoapHandler(w, r, service)
	})

	log.Println("SOAP server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
