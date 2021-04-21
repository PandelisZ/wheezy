package products

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type ProductService struct {
	Client   http.Client
	products []product
}

type product struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Quantity int    `json:"quantity"`
	Category string `json:"category"`
}

const productAPI = "https://run.mocky.io/v3/87beb2be-15bb-4b42-99b5-d826daaf8689"

func (ps ProductService) fetch() ([]product, error) {

	// Used cached products if we have already fetched
	if ps.products != nil {
		return ps.products, nil
	}

	req, err := http.NewRequest(http.MethodGet, productAPI, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, doErr := ps.Client.Do(req)
	if doErr != nil {
		log.Print(err)
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Print(readErr)
		return nil, readErr
	}

	products := []product{}
	jsonErr := json.Unmarshal(body, &products)
	if jsonErr != nil {
		log.Print(jsonErr)
		return nil, jsonErr
	}

	ps.products = products

	return products, nil
}

func (ps ProductService) All() ([]product, error) {
	products, err := ps.fetch()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps ProductService) Handler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	prd, readErr := ps.All()
	if readErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, readErr.Error())
		return
	}

	payload, err := json.Marshal(&prd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}
