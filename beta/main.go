package main

import (
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	_ "log"
	"net/http"
	_ "strings"
)

type ConsulKV struct {
	CreateIndex  int
	ModifyIndex  int
	LockIndex    int
	Key          string
	DecodedValue string
	Value        int
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://consul:8500/v1/kv/postgres/host_count/username")
	if err != nil {
		fmt.Fprintf(w, "unable to read the consul web interface: %v", err)
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		fmt.Fprintf(w, "no key found for %s", "postgres/host_count/username")
		return
	}
	/*
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(w, "unable to read the consul response: %v", err)
			return
		}
		str := string(b)
	*/
	dec := json.NewDecoder(resp.Body)
	var key ConsulKV
	if err := dec.Decode(&key); err != nil {
		fmt.Fprintf(w, "unable to decode the json response from consul: %v", err)
		return
	}
	fmt.Fprintf(w, "result:%#v", key)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("/static")))
	http.ListenAndServe(":80", nil)
}
