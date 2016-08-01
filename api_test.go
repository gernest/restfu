package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi(t *testing.T) {
	a := NewApi(NewStore()).Server()
	sample := []struct {
		name, desc string
	}{
		{"super", "best supre max"},
		{"supre max", "super maximum "},
	}

	var tempProds []*Product

	// Create new product
	for i := range sample {
		w := httptest.NewRecorder()
		data, _ := json.Marshal(&Product{
			Name: sample[i].name,
			Desc: sample[i].desc,
		})
		req, _ := http.NewRequest("POST", "/products", bytes.NewReader(data))
		a.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Error("expecked 5d got 5d", http.StatusOK, w.Code)
		}
		respProd := &Product{}
		err := json.Unmarshal([]byte(w.Body.String()), respProd)
		if err != nil {
			t.Fatal(err)
		}
		tempProds = append(tempProds, respProd)
	}

	// Get Product by id
	for i := range tempProds {
		w := httptest.NewRecorder()
		p := fmt.Sprintf("/products/%d.json", tempProds[i].ID)
		req, _ := http.NewRequest("GET", p, nil)
		a.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expecked %d got 5d", http.StatusOK, w.Code)
		}
		respProd := &Product{}
		err := json.Unmarshal([]byte(w.Body.String()), respProd)
		if err != nil {
			t.Fatal(err, p, w.Body)
		}
		if respProd.ID != tempProds[i].ID {
			t.Errorf("expected %d got 5d", respProd.ID, tempProds[i].ID)
		}
	}

	// Find all products
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products.json", nil)
		a.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expecked %d got 5d", http.StatusOK, w.Code)
		}
		respProd := make([]*Product, 0)
		err := json.Unmarshal([]byte(w.Body.String()), &respProd)
		if err != nil {
			t.Fatal(err)
		}
		if len(respProd) != 2 {
			t.Errorf("expected two products got %d", len(respProd))
		}
	}
	{
		w := httptest.NewRecorder()
		n := &Product{
			Name: "bob",
			Desc: "Almighty",
		}
		data, _ := json.Marshal(n)
		req, _ := http.NewRequest("POST", "/products/1", bytes.NewReader(data))
		a.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expecked %d got %d", http.StatusOK, w.Code)
		}
		respProd := &Product{}
		err := json.Unmarshal([]byte(w.Body.String()), respProd)
		if err != nil {
			t.Fatal(err)
		}
		if respProd.Name != n.Name {
			t.Errorf("expected %s got  %s", n.Name, respProd.Name)
		}
		if respProd.ID != 1 {
			t.Errorf("expected 1 got %d", respProd.ID)
		}
	}
}
