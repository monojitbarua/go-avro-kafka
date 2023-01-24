package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var configurationFile = "./config.json"
var Conf map[string]string = GetConfiguration()

func GetConfiguration() map[string]string {
	raw, err := ioutil.ReadFile(configurationFile)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var config map[string]string
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		panic(err)
	}

	fmt.Printf("\n Kafka configuration read from file %q \n%v\n", configurationFile, config)

	return config
}
