package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type asnInfo struct {
	Announced     bool   `json:"announced"`
	AsCountryCode string `json:"as_country_code"`
	AsDescription string `json:"as_description"`
	AsNumber      int    `json:"as_number"`
	FirstIP       string `json:"first_ip"`
	LastIP        string `json:"last_ip"`
}

func getASNInfo(ip string) (*asnInfo, error) {
	response, err := http.Get(fmt.Sprintf("https://api.iptoasn.com/v1/as/ip/%s", ip))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	info := asnInfo{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
