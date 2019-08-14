package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func LoadFixture(testName string, fixture interface{}) (err error) {
	fixtureName := fmt.Sprintf("testdata/%s.json", testName)
	jsonFile, err := os.Open(fixtureName)
	if err != nil {
		return
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, fixture)
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	return
}
