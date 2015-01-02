package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/igneous-systems/beta/shared"
)

const (
	USERPROP        = "postgres/host_count/username"
	PWDPROP         = "postgres/host_count/password"
	KV_API_ENDPOINT = "http://consul:8500/v1/kv"
)

type ConsulKV struct {
	CreateIndex  int
	ModifyIndex  int
	LockIndex    int
	Key          string
	DecodedValue string
	Value        string
}

func readConsul(key string) (*ConsulKV, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", KV_API_ENDPOINT, key))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	str, _ := ioutil.ReadAll(resp.Body)
	dec := json.NewDecoder(strings.NewReader(string(str)))
	defer resp.Body.Close()
	var kv []*ConsulKV
	if err := dec.Decode(&kv); err != nil {
		return nil, err
	}
	if len(kv) == 0 {
		return nil, nil //same as not found
	}
	for _, k := range kv {
		data, err := base64.StdEncoding.DecodeString(k.Value)
		if err != nil {
			return nil, err
		}
		k.DecodedValue = string(data)
	}
	if len(kv) > 1 {
		log.Printf("Ignoring extra keys found! (%d total found)", len(kv))
	}
	return kv[0], nil
}

func writeConsul(key string, newVal string) error {
	client := &http.Client{}

	req, err := http.NewRequest("PUT",
		fmt.Sprintf("%s/%s", KV_API_ENDPOINT, key), strings.NewReader(newVal))
	if err != nil {
		log.Printf("problem create request: %+v", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("problem writing to consul: %+v", err)
		return err
	}
	all, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		log.Printf("error in consul write: %+v, %+v, %s", resp, err, string(all))
		return fmt.Errorf("bad status: %d, %s", resp.StatusCode, string(all))
	}
	return nil
}

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

	if err := writeConsul(USERPROP, payload.Username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := writeConsul(PWDPROP, payload.Password); err != nil {
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
	resp, err := readConsul(USERPROP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp == nil {
		http.NotFound(w, r)
		return
	}
	result.Username = resp.DecodedValue

	resp, err = readConsul(PWDPROP)
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
