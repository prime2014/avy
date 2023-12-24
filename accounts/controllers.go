package accounts

import (
	"encoding/json"
	"fmt"
	"strings"

	"net/http"
)

func SignupController(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not Allowed!", http.StatusMethodNotAllowed)
	} else {
		var users Users
		err := json.NewDecoder(r.Body).Decode(&users)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = users.Validate()

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			errValues := strings.Split(err.Error(), ";")
			myerror := map[string][]string{"message": errValues}
			bytes, _ := json.Marshal(&myerror)
			http.Error(w, string(bytes), http.StatusBadRequest)
			return
		}
		users = users.Save()

		bytes, err := json.Marshal(users)

		if err != nil {
			http.Error(w, "Error encoding data", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, string(bytes))

	}
}

func LoginController(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
	} else {
		var users Users
		err := json.NewDecoder(r.Body).Decode(&users)

		fmt.Println(users)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		token, err := users.Authenticate()

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			errorMessage := map[string]string{"message": err.Error()}
			myerror, _ := json.Marshal(errorMessage)
			http.Error(w, string(myerror), http.StatusBadRequest)
			return
		}

		resp, err := json.Marshal(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(resp))
	}
}
