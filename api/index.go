package api

import "net/http"

func Index(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, World!"))
}

func init() {
	http.HandleFunc("/", Index)
}
