package config

import (
	"io/ioutil"
	"github.com/hjson/hjson-go"
	"strings"
	"fmt"
)


func Get_hosts() []string{

	dat, _ := ioutil.ReadFile("conf.json")

	var resp map[string]interface{}

	hjson.Unmarshal( dat, &resp )

	myHost, _ := resp["host"].(string)

	retHost := strings.Split(myHost, ",")


	var retHostTwo []string

	for _, host := range retHost {

		retHostTwo = append(retHostTwo, strings.Replace(host, " ", "", -1))
	}

	fmt.Println("array that is getting returned for hosts", retHostTwo)

	return retHostTwo
}

func Get_config(filename string, config_key string) string {

	dat, _ := ioutil.ReadFile( filename + ".json")

	var resp map[string]interface{}

	hjson.Unmarshal( dat, &resp)

	retVal, _ := resp[config_key].(string)

	return retVal
}


