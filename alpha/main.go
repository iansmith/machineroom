package main

import (
	"fmt"
	"log"
	"net/http"
)

type HostCount struct {
	Host  string
	Count int
}

func handler(w http.ResponseWriter, r *http.Request) {
	//try to contact the DB

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":80", nil)
}
