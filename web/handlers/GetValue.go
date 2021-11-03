package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/PrayagS/kv-store/pkg/kvstore"
)

type IncomingGETRequest struct {
	Key string `json:"Key"`
}

type OutgoingPOSTRequest struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func GetValue(kvstore *kvstore.KVStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IncomingGETRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		v, ok := kvstore.Get(req.Key)
		if ok != nil {
			http.Error(w, "The given key is not present in the data store.", http.StatusInternalServerError)
			return
		}

		var res OutgoingPOSTRequest
		res.Key = req.Key
		res.Value = v

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	}
}

func GetAllValues(kvstore *kvstore.KVStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(kvstore.Pairs)
	}
}
