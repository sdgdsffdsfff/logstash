package master

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

var (
	Err_NodeInitFailed = errors.New("Node Init Failed")
)

type config struct {
	Label  string `json:"label"`
	Addr   string `json:"addr"`
	System struct {
		Sampling int `json:"sampling"`
	} `json:"system"`
	Port struct {
		Nums     []string `json:"nums"`
		Sampling int      `json:"sampling"`
	} `json:"port"`
}

func (c config) String() string {
	j, _ := json.Marshal(c)
	return string(j)
}

var node_config config

func InitConfig() {
	cfg, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(cfg))

	jerr := json.Unmarshal(cfg, &node_config)
	if jerr != nil {
		panic(jerr)
	}

	fmt.Println(node_config)
}
