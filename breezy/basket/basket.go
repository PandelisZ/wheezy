package basket

import (
	"encoding/json"
	"github/pandelisz/wheezy/breezy/products"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type BasketService struct {
	Client         *http.Client
	ProductService *products.ProductService
	baskets        []basket
}

type upstreamProduct struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}

type basketUpstream struct {
	Id       int               `json:"id"`
	Products []upstreamProduct `json:"products"`
}

type productInBasket struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}

type basket struct {
	Id       int `json:"id"`
	Products []products.Product
	Total    float32
}

const basketAPI = "https://run.mocky.io/v3/6b59c8b9-deb9-49b7-b9ee-bc591a5ecf6d"

func (bs BasketService) fetch() ([]basket, error) {

	// Used cached products if we have already fetched
	if bs.baskets != nil {
		return bs.baskets, nil
	}

	req, err := http.NewRequest(http.MethodGet, basketAPI, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, doErr := bs.Client.Do(req)
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

	basketDirty := []basketUpstream{}
	jsonErr := json.Unmarshal(body, &basketDirty)
	if jsonErr != nil {
		log.Print(jsonErr)
		return nil, jsonErr
	}

	//Data enrichment step with info from product service
	productsById, productError := bs.ProductService.ProductByID()

	if productError != nil {
		return nil, productError
	}

	log.Println(basketDirty)

	basketEnriched := make([]basket, len(basketDirty))
	for k, val := range basketDirty {

		var enrichedProducts []products.Product
		runningTotal := float32(0)
		for _, p := range val.Products {
			log.Println(p)
			prd := productsById[p.Id]
			prd.Quantity = p.Quantity
			log.Println(prd)
			runningTotal += prd.Price * float32(p.Quantity)
			enrichedProducts = append(enrichedProducts, prd)
		}

		basketEnriched[k] = basket{
			Id:       val.Id,
			Products: enrichedProducts,
			Total:    runningTotal,
		}

	}

	bs.baskets = basketEnriched

	return basketEnriched, nil
}

func (bs BasketService) All() ([]basket, error) {
	products, err := bs.fetch()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (bs BasketService) Handler(w http.ResponseWriter, r *http.Request) {

	prd, readErr := bs.All()
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
