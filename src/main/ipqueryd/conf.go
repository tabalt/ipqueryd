package main

import (
	"encoding/json"
	"io/ioutil"
)

type Conf struct {
	PidFile string `json:"pid_file"`

	HttpServerPort string `json:"http_server_port"`
	GrpcServerPort string `json:"grpc_server_port"`
	DataFile       string `json:"data_file"`
}

// parse config for ipqueryd
func parseIpquerydConf(file string) (*Conf, error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := &Conf{}
	if err = json.Unmarshal(text, c); err != nil {
		return nil, err
	}

	return c, nil
}
