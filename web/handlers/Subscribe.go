package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/PrayagS/kv-store/pkg/kvstore"
	"github.com/gorilla/websocket"
	kafkaGo "github.com/segmentio/kafka-go"
)

type KafkaRecord struct {
	Topic     string
	Partition int
	Offset    int
	Key       string
	Value     string
}

var upgrader = websocket.Upgrader{} // use default options

func Subscribe(kvstore *kvstore.KVStore, kafkaReader *kafkaGo.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		// go sendKafkaRecords(c, kafkaReader, w, r)
		// The event loop
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Error during message reading:", err)
				break
			}
			log.Printf("Received: %s", message)

			m, err := kafkaReader.ReadMessage(context.Background())
			if err != nil {
				log.Fatalln(err)
			}
			var msg KafkaRecord
			msg.Topic = m.Topic
			msg.Partition = m.Partition
			msg.Offset = int(m.Offset)
			msg.Key = string(m.Key)
			msg.Value = string(m.Value)

			err = c.WriteJSON(msg)
			if err != nil {
				log.Println("Error during message writing:", err)
				break
			}
		}
	}
}
