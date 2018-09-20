package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Cat incoming!")
		w.WriteHeader(200)
		html := `
<html>
<img src="https://api.thecatapi.com/v1/images/search?format=src&size=%s&formats=%s">
<footer>Powered by <a href="https://thecatapi.com">The Cat API</a></footer>
</html>
`
		s := os.Getenv("CAT_API_SIZE")
		fs := os.Getenv("CAT_API_FORMATS")
		fmt.Fprintf(w, fmt.Sprintf(html, s, fs))
	})

	http.ListenAndServe("0.0.0.0:8080", nil)
}
