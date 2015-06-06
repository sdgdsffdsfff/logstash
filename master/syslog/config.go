package master

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Mail struct {
		Receivers []string `json:"receivers"`
		Smtp      string   `json:"smtp"`
		Username  string   `json:"username"`
		Password  string   `json:"password"`
	} `json:"mail"`
	Sms struct {
		Receivers []string `json:"receivers"`
	} `json:"sms"`
	Mongo struct {
		Url string `json:"url"`
	}
}

var Config config

func InitConfigFormJson(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &Config); err != nil {
		panic(err)
	}
}
