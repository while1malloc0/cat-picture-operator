package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Cat incoming!")

		catResp, err := http.Get("http://thecatapi.com/api/images/get")

		if err != nil {
			fmt.Fprintf(w, "Could not get cat: %v", err)
		}

		w.WriteHeader(catResp.StatusCode)
		io.Copy(w, catResp.Body)
	})

	http.ListenAndServe("0.0.0.0:8080", nil)
}
