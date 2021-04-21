package products

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ProductService struct {
	Client   http.Client
	products []product
}

type productUpstream struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Quantity int    `json:"quantity"`
	Category string `json:"category"`
}

type product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
}

type category struct {
	Name  string `json:"name"`
	Items int    `json:"items"`
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

	productsDirty := []productUpstream{}
	jsonErr := json.Unmarshal(body, &productsDirty)
	if jsonErr != nil {
		log.Print(jsonErr)
		return nil, jsonErr
	}

	products := make([]product, len(productsDirty))
	for k, val := range productsDirty {
		floatPrice, _ := strconv.ParseFloat(val.Price, 32)
		products[k] = product{
			Id:       val.Id,
			Name:     val.Name,
			Price:    float32(floatPrice),
			Quantity: val.Quantity,
			Category: val.Category,
		}
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

func (ps ProductService) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	cat, readErr := ps.categories()
	if readErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, readErr.Error())
		return
	}

	payload, err := json.Marshal(&cat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (ps ProductService) categories() ([]*category, error) {
	products, err := ps.fetch()
	if err != nil {
		return nil, err
	}
	categories := make(map[string]*category)
	for _, p := range products {
		if val, ok := categories[p.Category]; ok {
			val.Items = val.Items + 1
			categories[p.Category] = val
		} else {
			categories[p.Category] = &category{
				Name:  p.Category,
				Items: 1,
			}
		}

	}

	var categoryArray []*category
	for _, val := range categories {
		categoryArray = append(categoryArray, val)
	}
	return categoryArray, nil
}
