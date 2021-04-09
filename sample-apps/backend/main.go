/*
Copyright (c) 2019 Keitaro AB

Use of this source code is governed by an MIT license
that can be found in the LICENSE file or at
https://opensource.org/licenses/MIT.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("API_URL")
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)
	fmt.Fprint(w, "<h1>"+responseString+"</h1>\n")
}

func main() {
	fmt.Println("Backend waiting for connections on :8080")
	http.HandleFunc("/", serveHTTP)
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}

