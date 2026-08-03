package main

import (
	"log"
	"net/http"
	"os"

	appcalculation "calculator-app/backend/internal/application/calculation"
	"calculator-app/backend/internal/interfaces/httpapi"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	service := appcalculation.NewService()
	log.Printf("calculator API listening on :%s", port)
	if err := http.ListenAndServe(":"+port, httpapi.NewHandler(service)); err != nil { log.Fatal(err) }
}
