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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

type SetValuePOSTRequest struct {
	Key   string
	Value string
}

// var key string
var value string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Performs the set operation and set the value for the given key.",
	Long: `Perform a SET operation.
	When no flag is passed, all key-value pairs are fetched.
	Given a key-value pair, the concerned key-value pair is either added or updated if already exists.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		payload := SetValuePOSTRequest{
			Key:   key,
			Value: value,
		}
		p, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Invalid key input")
		}
		request, err := http.Post("http://"+serverAddress+"/set", "application/json", bytes.NewBuffer(p))
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer request.Body.Close()
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(body))
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&key, "key", "k", "", "The key to search for in the key-value store.")
	setCmd.Flags().StringVarP(&value, "value", "v", "", "The value to set for the given key.")
	_ = setCmd.MarkFlagRequired("key")
	_ = setCmd.MarkFlagRequired("value")
}
