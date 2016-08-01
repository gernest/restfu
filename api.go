package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Currency string

type Product struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Desc  string   `json:desc"`
	Tags  []string `json:"tags"`
	Price []Price  `json:"price"`
}

type Price struct {
	Amount   float64  `json:"amount"`
	Currency Currency `json:"currency"`
}

type Store struct {
	products map[int]*Product
	mu       sync.RWMutex
}

func NewStore() *Store {
	return &Store{products: make(map[int]*Product)}
}

func (s *Store) Create(p *Product) *Product {
	s.mu.Lock()
	p.ID = len(s.products) + 1
	s.products[p.ID] = p
	s.mu.Unlock()
	return p
}

func (s *Store) All() []*Product {
	var rst []*Product
	s.mu.RLock()
	for i := range s.products {
		rst = append(rst, s.products[i])
	}
	s.mu.RUnlock()
	return rst
}

func (s *Store) FindByID(id int) *Product {
	s.mu.RLock()
	p := s.products[id]
	s.mu.RUnlock()
	return p
}

func (s *Store) Update(p *Product) *Product {
	s.mu.Lock()
	s.products[p.ID] = p
	s.mu.Unlock()
	return p
}

type Api struct {
	store *Store
}

func NewApi(s *Store) *Api {
	return &Api{store: s}
}

//NewProduct creates a new product and serves the details of the created product
//as the response.
func (a *Api) NewProduct(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)
	o := &Product{}
	_ = json.NewDecoder(r.Body).Decode(o)
	e.Encode(a.store.Create(o))
}

//UpdateProduct updates product details.
func (a *Api) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)
	vars := mux.Vars(r)
	sid := vars["id"]
	id, _ := strconv.Atoi(sid)
	o := &Product{}
	_ = json.NewDecoder(r.Body).Decode(o)
	o.ID = id
	e.Encode(a.store.Update(o))
}

//ViewProduct  serves the product details. The product is supposed to be
//identified by id.
func (a *Api) ViewProduct(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)
	vars := mux.Vars(r)
	sid := vars["id"]
	id, _ := strconv.Atoi(sid)
	_ = e.Encode(a.store.FindByID(id))
}

//ListProducts this shows the list of all products as a json array. The products
//are not ordered due to the fact that the implementation on how products are
//stored is a map, which has unorderd key.
//
// However it is possible to order the result by implementing an additional
// product list struct which satisfies the sort interface.
func (a *Api) ListProducts(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)
	o := a.store.All()
	_ = e.Encode(o)
}

//Server returns http.handler with all the endpoints registered. It uses gorilla
//mux for the handler.
func (a *Api) Server() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/products", a.NewProduct).Methods("POST")
	r.HandleFunc("/products/{id}", a.UpdateProduct).Methods("POST")
	r.HandleFunc("/products.json", a.ListProducts).Methods("GET")
	r.HandleFunc("/products/{id}.json", a.ViewProduct).Methods("GET")
	return r
}
