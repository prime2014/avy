package cart

import (
	// "fmt"
	"accounts"
	"fmt"
	"strconv"

	"encoding/json"
	"net/http"
	"strings"
)

func CartPostController(w http.ResponseWriter, r *http.Request) {

	var token, _ = r.Context().Value("token").(string)
	// fmt.Printf("%v", token)
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Unauthorized!", http.StatusUnauthorized)
		return
	}

	tokenString := strings.Split(token, " ")
	// fmt.Println(tokenString[1])
	id, _ := accounts.ParseJWTToken(tokenString[1])

	myid, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var crt CartRequest

	err = json.NewDecoder(r.Body).Decode(&crt)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = crt.Validate()

	if err != nil {
		fmt.Printf("THIS ERROR IS WORSE: %v", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	mycart, err := crt.Save(myid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&mycart)

}
