// jsonPlayer
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config map[string]interface{}

//func (c *Config) GetConfig(string name) interface{} {
//	return c[name]
//}

// we will create a function called LoadConfig that will load up our JSON
// configuration file

func LoadConfig(path string) (map[string]interface{}, error) {
	var m map[string]interface{}
	data, err := ioutil.ReadFile(path) //data is byte[]
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(data, &m)
	return m, err
}

func loadConfigObj(path string) (Config, error) {
	var c Config
	data, err := os.Open(path) //data is file implemented io.Reader
	if err != nil {
		return c, err
	}
	err = json.NewDecoder(data).Decode(&c) //c is any type match the json
	fmt.Println(c["name"].(string))
	return c, err
}

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(config)

	config2, err2 := loadConfigObj("config.json")
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(config2)


}
