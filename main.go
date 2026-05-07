package main

import "net/http"

func main() {
	servMux := http.ServeMux{}
	server := http.Server{Handler: &servMux, Addr: ":8080"}
	servMux.Handle("/", http.FileServer(http.Dir(".")))
	server.ListenAndServe()
}
