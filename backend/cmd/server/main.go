package main

import (
	"log"
	"net/http"
	"os"

	"calculator-app/backend/internal/httpapi"
)

func main() { port := os.Getenv("PORT"); if port == "" { port = "8080" }; log.Printf("calculator API listening on :%s", port); if err := http.ListenAndServe(":"+port, httpapi.NewHandler()); err != nil { log.Fatal(err) } }
