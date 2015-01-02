package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/igneous-systems/beta/shared"
	"github.com/igneous-systems/lib/consul" //shared with alpha, server only
)

const (
	USERPROP = "postgres/host_count/username"
	PWDPROP  = "postgres/host_count/password"
)

func post(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST API call %+v", r)

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var payload shared.ApiPayload

	if err := dec.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload.Username = strings.TrimSpace(payload.Username)
	payload.Password = strings.TrimSpace(payload.Password) //dangerous

	if err := consul.WriteSimpleValue(USERPROP, payload.Username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := consul.WriteSimpleValue(PWDPROP, payload.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//we return 200 from here, use {} to make the jquery api happy
	fmt.Fprint(w, "{}")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		post(w, r)
	} else {
		get(w, r)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET API call %+v", r)

	var result shared.ApiPayload
	resp, err := consul.ReadKV(USERPROP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp == nil {
		http.NotFound(w, r)
		return
	}
	result.Username = resp.DecodedValue

	resp, err = consul.ReadKV(PWDPROP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp == nil {
		result.Password = ""
	} else {
		result.Password = resp.DecodedValue
	}

	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	if err := enc.Encode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("encoding result for client: %s", buff.String())
	w.Write(buff.Bytes())
}

type staticFiles struct {
}

func (s staticFiles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("/static")).ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/beta/api", apiHandler)
	http.Handle("/beta/", http.StripPrefix("/beta/", staticFiles{}))
	http.ListenAndServe(":80", nil)
}
