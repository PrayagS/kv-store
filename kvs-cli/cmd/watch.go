/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var done = make(chan interface{})       // Channel to indicate that the receiverHandler is done
var interrupt = make(chan os.Signal, 1) // Channel to listen for interrupt signal to terminate gracefully

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Fatalln("Error in receive:", err)
			return
		}
		// fmt.Printf("Received: %s\n", msg)
		prettyJSON, err := PrettyString(string(msg))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(prettyJSON)
	}
}

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch the changes happening to all keys in the store",
	Long: `Watch the server for changes done by other clients.
	The client's address and the edited key-value pair is shown in JSON as received.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

		socketUrl := "ws://" + serverAddress + "/subscribe"
		conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
		if err != nil {
			log.Fatal("Error connecting to Websocket Server:", err)
		}
		defer conn.Close()
		go receiveHandler(conn)

		// Our main loop for the client
		// We send our relevant packets here
		for {
			select {
			case <-time.After(time.Duration(1) * time.Millisecond * 1000):
				// Send an echo packet every second
				err := conn.WriteMessage(websocket.TextMessage, []byte("Keepalive ping"))
				if err != nil {
					log.Fatalln("Error during writing to websocket:", err)
					return
				}
			case <-interrupt:
				// We received a SIGINT (Ctrl + C). Terminate gracefully...
				log.Println("Received SIGINT interrupt signal. Closing all pending connections")

				// Close our websocket connection
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Fatalln("Error during closing websocket:", err)
					return
				}

				select {
				case <-done:
					log.Println("Receiver Channel Closed! Exiting....")
				case <-time.After(time.Duration(1) * time.Second):
					log.Println("Timeout in closing receiving channel. Exiting....")
				}
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
