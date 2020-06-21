package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "healthy")
	})

	fmt.Println("Starting server...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Printf("Error staring server: %s", err)
	}
}
