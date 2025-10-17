package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	name := os.Getenv("APP_NAME")
	if name == "" {
		name = "instance"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from %s on port %s!\n", name, port)
	})

	fmt.Println("Listening on port:", port)
	http.ListenAndServe(":"+port, nil)
}
