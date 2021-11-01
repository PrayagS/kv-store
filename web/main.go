package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PrayagS/kv-store/pkg/kafka"
	"github.com/PrayagS/kv-store/pkg/kvstore"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	kafkaGo "github.com/segmentio/kafka-go"
)

type RequestMessage struct {
	Key string `json:"Key"`
}

type ResponseMessage struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func setValue(kvstore *kvstore.KVStore, kafkaWriter *kafkaGo.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ResponseMessage
		err := json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = kafka.AppendCommandLog(r.Context(), kafkaWriter, []byte(fmt.Sprintf("address-%s", r.RemoteAddr)), []byte(req.Key+" "+req.Value))
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

// func producerHandler(kafkaWriter *kafkaGo.Writer) http.HandlerFunc {
// 	return func(wrt http.ResponseWriter, req *http.Request) {
// 		body, err := ioutil.ReadAll(req.Body)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		// msg := kafka.Message{
// 		// 	Key:   []byte(fmt.Sprintf("address-%s", req.RemoteAddr)),
// 		// 	Value: body,
// 		// }
// 		// err = kafkaWriter.WriteMessages(req.Context(), msg)
// 		err = kafka.AppendCommandLog(req.Context(), kafkaWriter, []byte(fmt.Sprintf("address-%s", req.RemoteAddr)), body)

// 		if err != nil {
// 			_, _ = wrt.Write([]byte(err.Error()))
// 			log.Fatalln(err)
// 		}
// 	}
// }

func getValue(kvstore *kvstore.KVStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestMessage
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

		var res ResponseMessage
		res.Key = req.Key
		res.Value = v

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	}
}

func main() {

	s := kvstore.New()
	s.Set("lmao", "420")
	s.Set("nice", "69")

	kafkaWriter := kafka.GetKafkaWriter()
	defer kafkaWriter.Close()

	kafkaReader := kafka.GetKafkaReader()
	defer kafkaReader.Close()

	r := mux.NewRouter()
	r.Path("/get").Handler(getValue(&s))
	r.Path("/set").Handler(setValue(&s, kafkaWriter))

	// Run the web server.
	fmt.Println("start producer-api ... !!")
	// log.Fatal(http.ListenAndServe(":8080", nil))
	// Apply the CORS middleware to our top-level router, with the defaults.
	log.Fatal(http.ListenAndServe("0.0.0.0:23333", gorillaHandlers.CORS()(r)))
}

// func main() {
// 	s := kvstore.New()
// 	s.Set("lmao", 420)
// 	s.Set("nice", 69)
// 	x, ok := s.Get("lmao")
// 	if ok == nil {
// 		fmt.Printf("%s", x)
// 	}
// 	s.Set("lmao", 42069)
// 	x, ok = s.Get("lmao")
// 	if ok == nil {
// 		fmt.Printf("%v", x)
// 	}
// 	y, ok := s.Get("gg")
// 	fmt.Printf("%v %s", y, ok)
// }
