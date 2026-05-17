package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("manual wiring")) })
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("redirect route")) })
	fmt.Println("step 01 ready on :8081")
	_ = http.ListenAndServe(":8081", mux)
}
