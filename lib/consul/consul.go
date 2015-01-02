package consul

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	KV_API_ENDPOINT = "http://machineroom.node.consul:8500/v1/kv"
)

type Key struct {
	CreateIndex  int
	ModifyIndex  int
	LockIndex    int
	Key          string
	DecodedValue string
	Value        string
}

func ReadKV(key string) (*Key, error) {
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
	var kv []*Key
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

func WriteSimpleValue(key string, newVal string) error {
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
