package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PrayagS/kv-store/pkg/kafka"
	"github.com/PrayagS/kv-store/pkg/kvstore"
	kafkaGo "github.com/segmentio/kafka-go"
)

type IncomingPOSTRequest struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func SetValue(kvstore *kvstore.KVStore, kafkaWriter *kafkaGo.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IncomingPOSTRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = kafka.AppendCommandLog(r.Context(), kafkaWriter, []byte(fmt.Sprintf("Client address=%s", r.RemoteAddr)), []byte(fmt.Sprintf("%s: %s", req.Key, req.Value)))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			log.Fatalln(err)
		}

		kvstore.Set(req.Key, req.Value)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode("SUCCESS")
	}
}
