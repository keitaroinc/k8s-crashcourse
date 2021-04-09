/*
Copyright (c) 2019 Keitaro AB

Use of this source code is governed by an MIT license
that can be found in the LICENSE file or at
https://opensource.org/licenses/MIT.
*/

package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Here at pod "+os.Getenv("HOSTNAME")+" its fine.\n")
	})

	fmt.Println("API wannabe started and listening on 8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}

