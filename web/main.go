package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PrayagS/kv-store/pkg/kafka"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	kafkaGo "github.com/segmentio/kafka-go"
)

func producerHandler(kafkaWriter *kafkaGo.Writer) http.HandlerFunc {
	return func(wrt http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatalln(err)
		}
		// msg := kafka.Message{
		// 	Key:   []byte(fmt.Sprintf("address-%s", req.RemoteAddr)),
		// 	Value: body,
		// }
		// err = kafkaWriter.WriteMessages(req.Context(), msg)
		err = kafka.AppendCommandLog(req.Context(), kafkaWriter, []byte(fmt.Sprintf("address-%s", req.RemoteAddr)), body)

		if err != nil {
			wrt.Write([]byte(err.Error()))
			log.Fatalln(err)
		}
	}
}

func main() {

	kafkaWriter := kafka.GetKafkaWriter()
	defer kafkaWriter.Close()

	r := mux.NewRouter()
	r.Path("/").Handler(producerHandler(kafkaWriter))

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
