package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type RecordMap map[string]string

type Records struct {
	DomainName string
	Map        RecordMap
}

func read(path string) ([]Records, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	yamlData := string(content)
	rootMap := make(map[string]map[string]map[string]string)
	err = yaml.Unmarshal([]byte(yamlData), &rootMap)
	if err != nil {
		return nil, err
	}
	firstMap := rootMap["records"]
	records := make([]Records, 0)
	for k, v := range firstMap {
		records = append(records, Records{
			DomainName: k,
			Map:        v,
		})
	}
	return records, nil
}
