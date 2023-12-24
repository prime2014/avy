package products

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gosimple/slug"
)

func ProductPostController(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		var product Products
		err := json.NewDecoder(r.Body).Decode(&product)

		product.Slug = slug.Make(product.Name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		myproduct := product.Save()

		bytes, err := json.Marshal(myproduct)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, string(bytes))
	default:
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

}
