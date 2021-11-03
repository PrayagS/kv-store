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

type GetValuePOSTRequest struct {
	Key string
}

// var key string
var isAll bool

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Performs the get operation and returns the value for the given key.",
	Long:  `TODO: Longer description of GET`,
	Run: func(cmd *cobra.Command, args []string) {
		if key == "" || isAll {
			request, err := http.Get("http://" + serverAddress + "/getall")
			if err != nil {
				log.Fatalf("Error contacting server: %v", err)
			}
			defer request.Body.Close()
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.Fatalf("Error reading server response: %v", err)
			}
			prettyJSON, err := PrettyString(string(body))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(prettyJSON)
		} else {
			payload := GetValuePOSTRequest{
				Key: key,
			}
			p, err := json.Marshal(payload)
			if err != nil {
				log.Fatalf("Invalid key input. Error: %v", err)
			}
			request, err := http.Post("http://"+serverAddress+"/get", "application/json", bytes.NewBuffer(p))
			if err != nil {
				log.Fatalf("Error contacting server: %v", err)
			}
			defer request.Body.Close()
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.Fatalf("Error reading server response: %v", err)
			}
			prettyJSON, err := PrettyString(string(body))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(prettyJSON)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().StringVarP(&key, "key", "k", "", "The key to search for in the key-value store.")
	getCmd.Flags().BoolVarP(&isAll, "all", "a", false, "Set to true to fetch all keys")
}
