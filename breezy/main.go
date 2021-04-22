package main

import (
	"github/pandelisz/wheezy/breezy/basket"
	"github/pandelisz/wheezy/breezy/products"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	c := &http.Client{
		Timeout: time.Second * 2,
	}

	pService := &products.ProductService{
		Client: c,
	}
	bService := &basket.BasketService{
		Client:         c,
		ProductService: pService,
	}

	r.HandleFunc("/products", pService.Handler).Methods(http.MethodGet)
	r.HandleFunc("/categories", pService.CategoryHandler).Methods(http.MethodGet)
	r.HandleFunc("/baskets", bService.Handler).Methods(http.MethodGet)
	http.Handle("/", r)

	handler := cors.Default().Handler(r)

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
