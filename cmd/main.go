package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PrayagS/kv-store/pkg/kafka"
	"github.com/PrayagS/kv-store/pkg/kvstore"
	"github.com/PrayagS/kv-store/web/handlers"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type RequestMessage struct {
	Key string `json:"Key"`
}

type ResponseMessage struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type KafkaRecord struct {
	Topic     string
	Partition int
	Offset    int
	Key       string
	Value     string
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
	r.Path("/get").Handler(handlers.GetValue(&s))
	r.Path("/getall").Handler(handlers.GetAllValues(&s))
	r.Path("/set").Handler(handlers.SetValue(&s, kafkaWriter))
	r.Path("/subscribe").Handler(handlers.Subscribe(&s, kafkaReader))

	// Run the web server.
	fmt.Println("start producer-api ... !!")
	// Apply the CORS middleware to our top-level router, with the defaults.
	log.Fatal(http.ListenAndServe("0.0.0.0:23333", gorillaHandlers.CORS()(r)))
}