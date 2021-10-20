package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/janbaer/github-oauth-bridge/api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	mux := http.NewServeMux()

	mux.HandleFunc("/authCallback", api.AuthCallback)
	mux.HandleFunc("/login", api.Login)

	fmt.Println("Running github-oauth-bridge on port 9001...")

	log.Fatal(http.ListenAndServe(":9001", mux))
}
