package main

import (
	"github/pandelisz/wheezy/breezy/products"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	c := http.Client{
		Timeout: time.Second * 2,
	}

	pService := products.ProductService{
		Client: c,
	}

	r.HandleFunc("/products", pService.Handler).Methods(http.MethodGet)
	http.Handle("/", r)

	log.Println("Listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
